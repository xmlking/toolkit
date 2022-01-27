package xds

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sync/atomic"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_api_v3_auth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/rs/zerolog/log"
)

// refresher struct
type fileRefresher struct {
	version         uint32
	refreshInterval time.Duration
	nodeID          string
	fs              fs.FS
	ctx             context.Context
	snapshotCache   cachev3.SnapshotCache
}

var _ Refresher = (*fileRefresher)(nil)

func (r *fileRefresher) GetCache() cachev3.Cache {
	return r.snapshotCache
}

func (r *fileRefresher) Start() (err error) {
	for _, v := range strSlice {
		log.Debug().Str("component", "xds").Msgf("v: %s", v)

		var clusterName = "service_bbc"
		var remoteHost = v

		log.Info().Str("component", "xds").Msgf(">>>>>>>>>>>>>>>>>>> creating cluster, remoteHost, nodeID %s,  %s, %s", clusterName, v, r.nodeID)

		hst := &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address:  remoteHost,
				Protocol: core.SocketAddress_TCP,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: uint32(443),
				},
			},
		}}
		uctx := &envoy_api_v3_auth.UpstreamTlsContext{}
		tctx, err := ptypes.MarshalAny(uctx)
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Send()
		}

		c := []types.Resource{
			&cluster.Cluster{
				Name:                 clusterName,
				ConnectTimeout:       ptypes.DurationProto(2 * time.Second),
				ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
				DnsLookupFamily:      cluster.Cluster_V4_ONLY,
				LbPolicy:             cluster.Cluster_ROUND_ROBIN,
				LoadAssignment: &endpoint.ClusterLoadAssignment{
					ClusterName: clusterName,
					Endpoints: []*endpoint.LocalityLbEndpoints{{
						LbEndpoints: []*endpoint.LbEndpoint{
							{
								HostIdentifier: &endpoint.LbEndpoint_Endpoint{
									Endpoint: &endpoint.Endpoint{
										Address: hst,
									}},
							},
						},
					}},
				},
				TransportSocket: &core.TransportSocket{
					Name: "envoy.transport_sockets.tls",
					ConfigType: &core.TransportSocket_TypedConfig{
						TypedConfig: tctx,
					},
				},
			},
		}

		// =================================================================================
		var listenerName = "listener_0"
		var targetHost = v
		var targetPrefix = "/"
		var virtualHostName = "local_service"
		var routeConfigName = "local_route"

		log.Info().Str("component", "xds").Msgf(">>>>>>>>>>>>>>>>>>> creating listener " + listenerName)

		rte := &route.RouteConfiguration{
			Name: routeConfigName,
			VirtualHosts: []*route.VirtualHost{{
				Name:    virtualHostName,
				Domains: []string{"*"},
				Routes: []*route.Route{{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: targetPrefix,
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: clusterName,
							},
							PrefixRewrite: "/robots.txt",
							HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
								HostRewriteLiteral: targetHost,
							},
						},
					},
				}},
			}},
		}

		manager := &hcm.HttpConnectionManager{
			CodecType:  hcm.HttpConnectionManager_AUTO,
			StatPrefix: "ingress_http",
			RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
				RouteConfig: rte,
			},
			HttpFilters: []*hcm.HttpFilter{{
				Name: wellknown.Router,
			}},
		}

		pbst, err := ptypes.MarshalAny(manager)
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Send()
		}

		priv, err := fs.ReadFile(r.fs, "config/certs/upstream-localhost-key.pem")
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Send()
		}
		pub, err := fs.ReadFile(r.fs, "config/certs/upstream-localhost-cert.pem")
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Send()
		}

		// use the following imports
		// envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
		// envoy_api_v2_auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
		// core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
		// auth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"

		// 1. send TLS certs filename back directly

		sdsTls := &envoy_api_v3_auth.DownstreamTlsContext{
			CommonTlsContext: &envoy_api_v3_auth.CommonTlsContext{
				TlsCertificates: []*envoy_api_v3_auth.TlsCertificate{{
					CertificateChain: &core.DataSource{
						Specifier: &core.DataSource_InlineBytes{InlineBytes: []byte(pub)},
					},
					PrivateKey: &core.DataSource{
						Specifier: &core.DataSource_InlineBytes{InlineBytes: []byte(priv)},
					},
				}},
			},
		}

		// or
		// 2. send TLS SDS Reference value
		// sdsTls := &envoy_api_v3_auth.DownstreamTlsContext{
		// 	CommonTlsContext: &envoy_api_v3_auth.CommonTlsContext{
		// 		TlsCertificateSdsSecretConfigs: []*envoy_api_v3_auth.SdsSecretConfig{{
		// 			Name: "server_cert",
		// 		}},
		// 	},
		// }

		// 3. SDS via ADS

		// sdsTls := &envoy_api_v3_auth.DownstreamTlsContext{
		// 	CommonTlsContext: &envoy_api_v3_auth.CommonTlsContext{
		// 		TlsCertificateSdsSecretConfigs: []*envoy_api_v3_auth.SdsSecretConfig{{
		// 			Name: "server_cert",
		// 			SdsConfig: &core.ConfigSource{
		// 				ConfigSourceSpecifier: &core.ConfigSource_Ads{
		// 					Ads: &core.AggregatedConfigSource{},
		// 				},
		// 				ResourceApiVersion: core.ApiVersion_V3,
		// 			},
		// 		}},
		// 	},
		// }

		scfg, err := ptypes.MarshalAny(sdsTls)
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Send()
		}

		var l = []types.Resource{
			&listener.Listener{
				Name: listenerName,
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Protocol: core.SocketAddress_TCP,
							Address:  localhost,
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: 10000,
							},
						},
					},
				},
				FilterChains: []*listener.FilterChain{{
					Filters: []*listener.Filter{{
						Name: wellknown.HTTPConnectionManager,
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: pbst,
						},
					}},
					TransportSocket: &core.TransportSocket{
						Name: "envoy.transport_sockets.tls",
						ConfigType: &core.TransportSocket_TypedConfig{
							TypedConfig: scfg,
						},
					},
				}},
			}}

		var secretName = "server_cert"

		log.Info().Str("component", "xds").Msgf(">>>>>>>>>>>>>>>>>>> creating Secret " + secretName)
		var s = []types.Resource{
			&envoy_api_v3_auth.Secret{
				Name: secretName,
				Type: &envoy_api_v3_auth.Secret_TlsCertificate{
					TlsCertificate: &envoy_api_v3_auth.TlsCertificate{
						CertificateChain: &core.DataSource{
							Specifier: &core.DataSource_InlineBytes{InlineBytes: []byte(pub)},
						},
						PrivateKey: &core.DataSource{
							Specifier: &core.DataSource_InlineBytes{InlineBytes: []byte(priv)},
						},
					},
				},
			},
		}

		// =================================================================================
		atomic.AddUint32(&r.version, 1)
		log.Info().Str("component", "xds").Msgf(">>>>>>>>>>>>>>>>>>> creating snapshot Version " + fmt.Sprint(r.version))

		snap, err := cachev3.NewSnapshot("v0", map[resource.Type][]types.Resource{
			//resource.EndpointType:        endpoints,
			resource.ClusterType: c,
			//resource.RouteType:           routes,
			//resource.ScopedRouteType:     scopedRoutes,
			resource.ListenerType: l,
			//resource.RuntimeType:         runtimes,
			resource.SecretType: s,
			//resource.ExtensionConfigType: extensions,
		})

		if err := snap.Consistent(); err != nil {
			log.Error().Str("component", "xds").Msgf("snapshot inconsistency: %+v\n%+v", snap, err)
			os.Exit(1)
		}
		err = r.snapshotCache.SetSnapshot(r.ctx, r.nodeID, snap)
		if err != nil {
			log.Fatal().Str("component", "xds").Msgf("Could not set snapshot %v", err)
		}

		time.Sleep(10 * time.Second)
	}
	return
}
