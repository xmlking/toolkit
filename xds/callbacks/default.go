package callbacks

import (
	"context"
	"sync"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type defaultCallbacks struct {
	//connected clients
	streamClients int
	//connected delta clients
	deltaClients int
	fetchReq     int
	fetchResp    int
	streamReq    int
	streamResp   int
	deltaReq     int
	deltaResp    int
	mu           sync.Mutex
}

var _ server.Callbacks = (*defaultCallbacks)(nil)

func NewDefaultCallbacks() server.Callbacks {
	return &defaultCallbacks{
		streamClients: 0,
		deltaClients:  0,
		fetchReq:      0,
		fetchResp:     0,
		streamReq:     0,
		streamResp:    0,
		deltaReq:      0,
		deltaResp:     0,
	}
}

func (cb *defaultCallbacks) Report() {
	log.Debug().Str("component", "xds").
		Dict("xds_metrics", zerolog.Dict().
			Int("stream_clients", cb.streamClients).
			Int("delta_clients", cb.deltaClients).
			Int("fetch_req", cb.fetchReq).
			Int("fetch_resp", cb.fetchResp).
			Int("stream_req", cb.streamReq).
			Int("stream_resp", cb.streamResp).
			Int("delta_req", cb.deltaReq).
			Int("delta_resp", cb.deltaResp),
		).
		Msg("xds metrics report")
}

func (cb *defaultCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.streamClients++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Str("type", typeURL).Msg("StreamOpen")
	cb.Report()
	return nil
}

func (cb *defaultCallbacks) OnStreamClosed(streamID int64) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.streamClients--
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Msg("StreamClosed")
	cb.Report()
}

func (cb *defaultCallbacks) OnDeltaStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaClients++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Str("type", typeURL).Msg("DeltaStreamOpen")
	cb.Report()
	return nil
}

func (cb *defaultCallbacks) OnDeltaStreamClosed(streamID int64) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaClients--
	log.Debug().Str("component", "xds").Int64("streamID", streamID).Msg("DeltaStreamClosed")
	cb.Report()
}

func (cb *defaultCallbacks) OnStreamRequest(streamID int64, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.streamReq++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Msg("StreamRequest")
	cb.Report()
	return nil
}

func (cb *defaultCallbacks) OnStreamResponse(ctx context.Context, streamID int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.streamResp++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Interface("response", resp).Msg("StreamResponse")
	cb.Report()
}

func (cb *defaultCallbacks) OnStreamDeltaRequest(streamID int64, req *discovery.DeltaDiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaReq++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Msg("StreamDeltaRequest")
	cb.Report()
	return nil
}

func (cb *defaultCallbacks) OnStreamDeltaResponse(streamID int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaResp++
	log.Debug().Str("component", "xds").Int64("streamID", streamID).
		Interface("request", req).Interface("response", resp).Msg("StreamDeltaResponse")
	cb.Report()
}

func (cb *defaultCallbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetchReq++
	log.Debug().Str("component", "xds").Interface("request", req).Msg("FetchRequest")
	cb.Report()
	return nil
}

func (cb *defaultCallbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetchResp++
	log.Debug().Str("component", "xds").Interface("request", req).
		Interface("response", resp).Msg("FetchResponse")
	cb.Report()
}
