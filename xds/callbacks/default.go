package callbacks

import (
	"context"
	"sync"

	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/rs/zerolog/log"
)

type defaultCallbacks struct {
	Fetches  int
	Requests int
	mu       sync.Mutex
}

var _ serverv3.Callbacks = (*defaultCallbacks)(nil)

func NewDefaultCallbacks() serverv3.Callbacks {
	return &defaultCallbacks{
		Fetches:  0,
		Requests: 0,
	}
}

func (cb *defaultCallbacks) OnDeltaStreamOpen(ctx context.Context, i int64, s string) error {
	panic("implement me")
}
func (cb *defaultCallbacks) OnDeltaStreamClosed(i int64) {
	panic("implement me")
}
func (cb *defaultCallbacks) OnStreamDeltaRequest(i int64, request *discoverygrpc.DeltaDiscoveryRequest) error {
	panic("implement me")
}
func (cb *defaultCallbacks) OnStreamDeltaResponse(i int64, request *discoverygrpc.DeltaDiscoveryRequest, response *discoverygrpc.DeltaDiscoveryResponse) {
	panic("implement me")
}
func (cb *defaultCallbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.Info().Fields(map[string]interface{}{"fetches": cb.Fetches, "requests": cb.Requests}).Msg("cb.Report()  callbacks")
}
func (cb *defaultCallbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	log.Info().Msgf("OnStreamOpen %d open for %s", id, typ)
	return nil
}
func (cb *defaultCallbacks) OnStreamClosed(id int64) {
	log.Info().Msgf("OnStreamClosed %d closed", id)
}
func (cb *defaultCallbacks) OnStreamRequest(id int64, r *discoverygrpc.DiscoveryRequest) error {
	log.Info().Msgf("OnStreamRequest %v", r.TypeUrl)
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	return nil
}
func (cb *defaultCallbacks) OnStreamResponse(int64, *discoverygrpc.DiscoveryRequest, *discoverygrpc.DiscoveryResponse) {
	log.Info().Msgf("OnStreamResponse...")
	cb.Report()
}
func (cb *defaultCallbacks) OnFetchRequest(ctx context.Context, req *discoverygrpc.DiscoveryRequest) error {
	log.Info().Msgf("OnFetchRequest...")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	return nil
}
func (cb *defaultCallbacks) OnFetchResponse(*discoverygrpc.DiscoveryRequest, *discoverygrpc.DiscoveryResponse) {
	log.Info().Msgf("OnFetchResponse...")
}
