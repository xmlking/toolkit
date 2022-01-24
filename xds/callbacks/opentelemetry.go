package callbacks

import (
	"context"
	"sync"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
)

const libraryName = "xds_controller"

// go-control-plane defines standard interface for  callback mechanism, which can be used to record and expose metrics out of the xDS requests.
// GRPC SoTW (State of The World) part of XDS server functions below.. (Rest functions are not implemented)

type otelCallbacks struct {
	//Stream counters
	activeGauge metric.Int64UpDownCounter
	reqCounter  metric.Int64Counter
	resCounter  metric.Int64Counter
	// mux for incrementing counters
	mu sync.Mutex
}

var _ serverv3.Callbacks = (*otelCallbacks)(nil)

func NewOTelCallbacks() (cb serverv3.Callbacks, err error) {
	cbo := &otelCallbacks{}
	meter := global.Meter("otel-xds-controller")

	if cbo.activeGauge, err = meter.NewInt64UpDownCounter(
		"active_streams",
		metric.WithDescription("Active grpc streams to xds-controller"),
		metric.WithUnit(unit.Dimensionless),
	); err != nil {
		return
	}

	if cbo.reqCounter, err = meter.NewInt64Counter(
		"stream_requests",
		metric.WithDescription("No.of requests via grpc streams to xds-controller"),
		metric.WithUnit(unit.Dimensionless),
	); err != nil {
		return
	}

	if cbo.resCounter, err = meter.NewInt64Counter(
		"stream_responses",
		metric.WithDescription("No.of Responses sent to clients by  xds-controller"),
		metric.WithUnit(unit.Dimensionless),
	); err != nil {
		return
	}

	return cbo, err
}

// OnStreamOpen is called once an xDS stream is open with a stream ID and the type URL (or "" for ADS).
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (cb *otelCallbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	(cb.activeGauge).Add(ctx, 1)
	log.Debug().Msgf("Callback: Stream open for  id: %d open for type: %s", id, typ)
	return nil
}

// OnStreamClosed is called immediately prior to closing an xDS stream with a stream ID.
func (cb *otelCallbacks) OnStreamClosed(id int64) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.activeGauge.Add(context.Background(), -1)
	log.Debug().Msgf("Callback: Stream Closed for  id: %d", id)
}

// OnStreamRequest is called once a request is received on a stream.
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (cb *otelCallbacks) OnStreamRequest(a int64, d *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.reqCounter.Add(context.Background(), 1)
	log.Debug().Msgf("Callback: Stream Request %v", d)
	return nil
}

// OnStreamResponse is called immediately prior to sending a response on a stream.
func (cb *otelCallbacks) OnStreamResponse(ctx context.Context, a int64, req *discovery.DiscoveryRequest, d *discovery.DiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.resCounter.Add(context.Background(), 1)
	log.Debug().Msgf("Callback: Stream Response: %v", d)
}

// OnFetchRequest Marker Impl: No expecting Rest Client
func (cb *otelCallbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	return nil
}

// OnFetchResponse Marker Impl: No expecting Rest Client
func (cb *otelCallbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
}

func (cb *otelCallbacks) OnDeltaStreamOpen(ctx context.Context, i int64, s string) error {
	panic("implement me")
}

func (cb *otelCallbacks) OnDeltaStreamClosed(i int64) {
	panic("implement me")
}

func (cb *otelCallbacks) OnStreamDeltaRequest(i int64, request *discovery.DeltaDiscoveryRequest) error {
	panic("implement me")
}

func (cb *otelCallbacks) OnStreamDeltaResponse(i int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {
	panic("implement me")
}
