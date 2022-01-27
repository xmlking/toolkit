package kube

import (
	"context"
	"fmt"
	"net"
	"strconv"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	routerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	managerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/xds"
	"github.com/xmlking/toolkit/xds/kube/apigateway"
	"google.golang.org/protobuf/types/known/anypb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	k8scache "k8s.io/client-go/tools/cache"
)

func (k *kubeRefresher) startServices(ctx context.Context) error {
	emit := func() {
		log.Warn().Str("component", "xds").Msg("emit before ready")
	}

	store := k8scache.NewUndeltaStore(func(v []interface{}) {
		emit()
	}, k8scache.DeletionHandlingMetaNamespaceKeyFunc)

	reflector := k8scache.NewReflector(&k8scache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return k.client.CoreV1().Services(k.namespace).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return k.client.CoreV1().Services(k.namespace).Watch(ctx, options)
		},
	}, &corev1.Service{}, store, k.refreshInterval)

	emit = func() {
		services := sliceToService(store.List())
		version := reflector.LastSyncResourceVersion()

		resources := kubeServicesToResources(services)
		apiGatewayResources := apigateway.FromKubeServices(services)
		merged := append(resources, apiGatewayResources...)

		snapshot, err := cachev3.NewSnapshot(version, xds.ResourcesToMap(merged))
		if err != nil {
			log.Fatal().Err(err).Str("component", "xds").Msg("Failed to create New Snapshot")
		}

		if err := snapshot.Consistent(); err != nil {
			log.Fatal().Err(err).Str("component", "xds").Msg("Snapshot inconsistent")
		}

		if err := k.servicesCache.SetSnapshot(ctx, k.nodeID, snapshot); err != nil {
			log.Fatal().Err(err).Str("component", "xds").Msg("Failed to Set Snapshot to cache")
		} else {
			if k.telemetry {
				k.snapshots.Add(ctx, 1)
			}
		}
	}

	reflector.Run(ctx.Done())
	return nil
}

func sliceToService(s []interface{}) []*corev1.Service {
	out := make([]*corev1.Service, len(s))
	for i, v := range s {
		out[i] = v.(*corev1.Service)
	}
	return out
}

// kubeServicesToResources convert list of Kubernetes services to
// - Listener for each ports
// - RouteConfiguration for those listeners
// - Cluster
func kubeServicesToResources(services []*corev1.Service) []types.Resource {
	var out []types.Resource

	router, _ := anypb.New(&routerv3.Router{})

	for _, svc := range services {
		fullName := fmt.Sprintf("%s.%s", svc.Name, svc.Namespace)
		for _, port := range svc.Spec.Ports {
			targetHostPort := net.JoinHostPort(fullName, port.Name)
			targetHostPortNumber := net.JoinHostPort(fullName, strconv.Itoa(int(port.Port)))
			routeConfig := &routev3.RouteConfiguration{
				Name: targetHostPortNumber,
				VirtualHosts: []*routev3.VirtualHost{
					{
						Name:    targetHostPort,
						Domains: []string{fullName, targetHostPort, targetHostPortNumber, svc.Name},
						Routes: []*routev3.Route{{
							Name: "default",
							Match: &routev3.RouteMatch{
								PathSpecifier: &routev3.RouteMatch_Prefix{},
							},
							Action: &routev3.Route_Route{
								Route: &routev3.RouteAction{
									ClusterSpecifier: &routev3.RouteAction_Cluster{
										Cluster: targetHostPort,
									},
								},
							},
						}},
					},
				},
			}

			manager, _ := anypb.New(&managerv3.HttpConnectionManager{
				HttpFilters: []*managerv3.HttpFilter{
					{
						Name: wellknown.Router,
						ConfigType: &managerv3.HttpFilter_TypedConfig{
							TypedConfig: router,
						},
					},
				},
				RouteSpecifier: &managerv3.HttpConnectionManager_RouteConfig{
					RouteConfig: routeConfig,
				},
			})

			svcListener := &listenerv3.Listener{
				Name: targetHostPortNumber,
				ApiListener: &listenerv3.ApiListener{
					ApiListener: manager,
				},
			}

			svcCluster := &clusterv3.Cluster{
				Name:                 targetHostPort,
				ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_EDS},
				LbPolicy:             clusterv3.Cluster_ROUND_ROBIN,
				EdsClusterConfig: &clusterv3.Cluster_EdsClusterConfig{
					EdsConfig: &corev3.ConfigSource{
						ConfigSourceSpecifier: &corev3.ConfigSource_Ads{
							Ads: &corev3.AggregatedConfigSource{},
						},
					},
				},
			}

			out = append(out, svcListener, routeConfig, svcCluster)
		}
	}

	return out
}
