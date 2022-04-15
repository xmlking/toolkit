package kube

import (
	"context"
	"fmt"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/xds/api"
	"google.golang.org/protobuf/types/known/wrapperspb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	k8scache "k8s.io/client-go/tools/cache"
)

func (k *kubeRefresher) startEndpoints(ctx context.Context) error {
	emit := func() {}

	store := k8scache.NewUndeltaStore(func(v []interface{}) {
		emit()
	}, k8scache.DeletionHandlingMetaNamespaceKeyFunc)

	reflector := k8scache.NewReflector(&k8scache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return k.client.CoreV1().Endpoints(k.namespace).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return k.client.CoreV1().Endpoints(k.namespace).Watch(ctx, options)
		},
	}, &corev1.Endpoints{}, store, k.refreshInterval)

	emit = func() {
		endpoints := kubeEndpointsToResources(sliceToEndpoints(store.List()))
		version := reflector.LastSyncResourceVersion()

		snapshot, err := cachev3.NewSnapshot(version, api.ResourcesToMap(endpoints))
		if err != nil {
			log.Fatal().Stack().Err(err).Str("component", "xds").Msg("Failed to create New Snapshot")
		}

		if err := snapshot.Consistent(); err != nil {
			println(api.DebugSnapshot(&snapshot))
			log.Fatal().Stack().Err(err).
				Interface("snapshot", snapshot).
				Str("component", "xds").Msg("Snapshot inconsistent")
		}

		if err := k.endpointsCache.SetSnapshot(ctx, k.nodeID, snapshot); err != nil {
			log.Fatal().Stack().Err(err).Str("component", "xds").Msg("Failed to Set Snapshot to cache")
		} else {
			if k.telemetry {
				k.snapshots.Add(ctx, 1)
			}
		}
	}

	reflector.Run(ctx.Done())
	return nil
}

func sliceToEndpoints(s []interface{}) []*corev1.Endpoints {
	out := make([]*corev1.Endpoints, len(s))
	for i, v := range s {
		out[i] = v.(*corev1.Endpoints)
	}
	return out
}

// kubeServicesToResources convert list of Kubernetes endpoints to Endpoint
func kubeEndpointsToResources(endpoints []*corev1.Endpoints) []types.Resource {
	var out []types.Resource

	for _, ep := range endpoints {
		for _, subset := range ep.Subsets {
			for _, port := range subset.Ports {
				var portName string
				if port.Name == "" {
					portName = fmt.Sprintf("%s.%s:%d", ep.Name, ep.Namespace, port.Port)
				} else {
					portName = fmt.Sprintf("%s.%s:%s", ep.Name, ep.Namespace, port.Name)
				}

				cla := &endpointv3.ClusterLoadAssignment{
					ClusterName: portName,
					Endpoints: []*endpointv3.LocalityLbEndpoints{
						{
							LoadBalancingWeight: wrapperspb.UInt32(1),
							Locality:            &corev3.Locality{},
							LbEndpoints:         []*endpointv3.LbEndpoint{},
						},
					},
				}
				out = append(out, cla)

				for _, addr := range subset.Addresses {
					hostname := addr.Hostname
					if hostname == "" && addr.TargetRef != nil {
						hostname = fmt.Sprintf("%s.%s", addr.TargetRef.Name, addr.TargetRef.Namespace)
					}
					if hostname == "" && addr.NodeName != nil {
						hostname = *addr.NodeName
					}

					cla.Endpoints[0].LbEndpoints = append(cla.Endpoints[0].LbEndpoints, &endpointv3.LbEndpoint{
						HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
							Endpoint: &endpointv3.Endpoint{
								Address: &corev3.Address{
									Address: &corev3.Address_SocketAddress{
										SocketAddress: &corev3.SocketAddress{
											Protocol: corev3.SocketAddress_TCP,
											Address:  addr.IP,
											PortSpecifier: &corev3.SocketAddress_PortValue{
												PortValue: uint32(port.Port),
											},
										},
									},
								},
								Hostname: hostname,
							},
						},
					})
				}
			}
		}
	}

	return out
}
