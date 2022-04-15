package callbacks

import (
	"context"
	"sync"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
)

type otelCallbacks struct {
	connectedClients metric.Int64UpDownCounter
	requestsTotal    metric.Int64Counter
	responsesTotal   metric.Int64Counter
	mu               sync.Mutex
}

var _ server.Callbacks = (*otelCallbacks)(nil)

func NewOTelCallbacks() server.Callbacks {
	meter := global.Meter("otel-xds-server")
	return &otelCallbacks{
		connectedClients: metric.Must(meter).NewInt64UpDownCounter(
			"connected_clients",
			metric.WithDescription("Number of clients currently connected to the xds-server"),
			metric.WithUnit(unit.Dimensionless),
		),
		requestsTotal: metric.Must(meter).NewInt64Counter(
			"requests_total",
			metric.WithDescription("Number of requests from clients to the xds-server"),
			metric.WithUnit(unit.Dimensionless),
		),
		responsesTotal: metric.Must(meter).NewInt64Counter(
			"responses_total",
			metric.WithDescription("Number of responses sent to clients by the xds-server"),
			metric.WithUnit(unit.Dimensionless),
		),
	}
}

func (cb *otelCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	// TODO: Implement Auth ? metadata, ok := metadata.FromIncomingContext(ctx);  metadata[credentialTokenHeaderKey];
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.connectedClients.Add(ctx, 1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Str("type", typeURL).Msg("StreamOpen")
	return nil
}

func (cb *otelCallbacks) OnStreamClosed(streamID int64) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.connectedClients.Add(context.Background(), -1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Msg("StreamClosed")
}

func (cb *otelCallbacks) OnDeltaStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.connectedClients.Add(ctx, 1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Str("type", typeURL).Msg("DeltaStreamOpen")
	return nil
}

func (cb *otelCallbacks) OnDeltaStreamClosed(streamID int64) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.connectedClients.Add(context.Background(), -1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Msg("DeltaStreamClosed")
}

func (cb *otelCallbacks) OnStreamRequest(streamID int64, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requestsTotal.Add(context.Background(), 1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Msg("StreamRequest")
	return nil
}

func (cb *otelCallbacks) OnStreamResponse(ctx context.Context, streamID int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.responsesTotal.Add(context.Background(), 1)
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Interface("response", resp).Msg("StreamResponse")
}

func (cb *otelCallbacks) OnStreamDeltaRequest(streamID int64, req *discovery.DeltaDiscoveryRequest) error {
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Msg("StreamDeltaRequest")
	return nil
}

func (cb *otelCallbacks) OnStreamDeltaResponse(streamID int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Interface("response", resp).Msg("StreamDeltaResponse")
}

func (cb *otelCallbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	log.Debug().Str("component", "xds").Interface("request", req).Msg("FetchRequest")
	return nil
}

func (cb *otelCallbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	log.Debug().Str("component", "xds").Interface("request", req).
		Interface("response", resp).Msg("FetchResponse")
}
