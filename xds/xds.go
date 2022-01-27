package xds

import (
	"context"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/xmlking/toolkit/xds/api"
	"github.com/xmlking/toolkit/xds/kube"
)

// NewRefresher is generic constructor
func NewRefresher(ctx context.Context, sType string, opts ...Option) (r api.Refresher) {
	// Set default options
	options := Options{
		RefreshInterval: 0,
		FileSys:         os.DirFS("."),
	}

	// copy options
	for _, o := range opts {
		o(&options)
	}

	switch sType {
	case KUBERNETES:
		r = kube.NewKubeRefresher(ctx, options.RefreshInterval, options.NodeID, options.Namespace, newXdsLogger())
	case DNS:
		r = NewDNSRefresher(ctx, options.RefreshInterval, options.NodeID, options.Hostnames, cache.NewSnapshotCache(true, cache.IDHash{}, newXdsLogger()))
	case FILE:
		r = &fileRefresher{
			version:         0,
			refreshInterval: options.RefreshInterval,
			fs:              options.FileSys,
			ctx:             ctx,
			nodeID:          options.NodeID,
			snapshotCache:   cache.NewSnapshotCache(true, cache.IDHash{}, newXdsLogger()),
		}
	default:
		r = &staticRefresher{
			version:         0,
			refreshInterval: options.RefreshInterval,
			fs:              options.FileSys,
			ctx:             ctx,
			nodeID:          options.NodeID,
			snapshotCache:   cache.NewSnapshotCache(true, cache.IDHash{}, newXdsLogger()),
		}
	}

	return
}
