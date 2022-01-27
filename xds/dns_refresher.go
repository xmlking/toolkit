package xds

import (
	"context"
	"net"
	"sync"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/xds/api"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type dnsRefresher struct {
	refreshInterval time.Duration
	ctx             context.Context
	nodeID          string

	hostnames     []string
	snapshotCache cache.SnapshotCache

	telemetry bool
	snapshots metric.Int64Counter
	errors    metric.Int64Counter
	mu        sync.Mutex
}

var _ api.Refresher = (*dnsRefresher)(nil)

func NewDNSRefresher(ctx context.Context, refreshInterval time.Duration, nodeID string, hostnames []string, snapshotCache cache.SnapshotCache) api.Refresher {
	meter := global.Meter("otel-k8s-refresher")
	dns := &dnsRefresher{
		ctx:             ctx,
		refreshInterval: refreshInterval,
		nodeID:          nodeID,
		hostnames:       hostnames,
		snapshotCache:   snapshotCache,
	}

	var err error
	dns.snapshots, err = meter.NewInt64Counter(
		"snapshots",
		metric.WithDescription("Number of snapshots generated by xds-controller"),
		metric.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		log.Warn().Err(err).Str("component", "xds").Msg("Telemetry not enabled")
	} else {
		dns.telemetry = true
	}

	dns.errors, err = meter.NewInt64Counter(
		"errors",
		metric.WithDescription("Number of errors happened in xds-controller"),
		metric.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		log.Warn().Err(err).Str("component", "xds").Msg("Telemetry not enabled")
	} else {
		dns.telemetry = true
	}

	return dns
}

func (d *dnsRefresher) GetCache() cache.Cache {
	return d.snapshotCache
}

func (d *dnsRefresher) Start() (err error) {
	for {
		select {
		case <-d.ctx.Done():
			// log situation
			switch d.ctx.Err() {
			case context.DeadlineExceeded:
				log.Debug().Str("component", "xds").Msg("Context timeout exceeded")
			case context.Canceled:
				log.Debug().Str("component", "xds").Msg("Context cancelled by interrupt signal")
			}
			log.Info().Str("component", "xds").Msg("Stopping dnsRefresher...")
			return
		default:
			if dnsEndpoints, err2 := d.getEndpoints(); err2 == nil {
				// TODO check if dnsEndpoints changed, skip next steps if not changed.
				version := uuid.New().String()
				endpoints := dnsEndpointsToResources(dnsEndpoints)

				var snapshot cache.Snapshot
				snapshot, err = cache.NewSnapshot(version, api.ResourcesToMap(endpoints))
				if err != nil {
					log.Error().Err(err).Str("component", "xds").Msg("Failed to create New Snapshot")
					return
				}
				if err = snapshot.Consistent(); err != nil {
					log.Error().Err(err).Str("component", "xds").Msg("Snapshot inconsistent")
					return
				}
				if err = d.snapshotCache.SetSnapshot(d.ctx, d.nodeID, snapshot); err != nil {
					log.Error().Err(err).Str("component", "xds").Msg("Failed to Set Snapshot to cache")
					return
				} else {
					if d.telemetry {
						d.snapshots.Add(d.ctx, 1)
					}
				}

			} else {
				if d.telemetry {
					d.errors.Add(d.ctx, 1)
				}
			}
			time.Sleep(d.refreshInterval)
		}
	}
}

// getEndpoints returns a slice of that host's IP addresses for each host
func (d *dnsRefresher) getEndpoints() (endpoints map[string][]string, err error) {
	endpoints = make(map[string][]string)
	for _, host := range d.hostnames {
		var addrs []string
		if addrs, err = d.getAddresses(host); err == nil {
			log.Debug().Str("component", "xds").Str("hostname", host).
				Strs("addresses", addrs).Msg("DNS resolved")
			endpoints[host] = addrs
		} else {
			log.Error().Err(err).Str("component", "xds").
				Msgf("Unable to resolve IP addresses for hostname: %s.", host)
		}
	}
	return
}

func (d *dnsRefresher) getAddresses(hostname string) (addrs []string, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var netResolver = net.DefaultResolver
	return netResolver.LookupHost(context.Background(), hostname)
}

func dnsEndpointsToResources(endpoints map[string][]string) (out []types.Resource) {
	for host, addrs := range endpoints {
		out = append(out, createEds(host, addrs))
	}
	return
}

func createEds(clusterName string, ips []string) types.Resource {
	var lbs []*endpoint.LbEndpoint
	// create lb for available service ips
	for _, ip := range ips {
		lbs = append(lbs, &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Address:  ip,
							Protocol: core.SocketAddress_TCP,
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: uint32(8080), // TODO configure port
							},
						},
					}},
				}},
			// TODO
			HealthStatus: core.HealthStatus_HEALTHY,
		})
	}
	// create an eds cluster with lbs
	return types.Resource(&endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{
			{
				Locality: &core.Locality{
					Region: "region",
					Zone:   "zone",
				},
				Priority:            0,
				LoadBalancingWeight: &wrapperspb.UInt32Value{Value: uint32(1000)},
				LbEndpoints:         lbs,
			},
		},
	})
}
