package xds

import (
	"context"
	"os"

	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

type Refresher interface {
	GetSnapshotCache() cachev3.SnapshotCache
	Start() error
}

// NewRefresher is generic constructor
func NewRefresher(ctx context.Context, opts ...Option) (r Refresher) {
	// Set default options
	options := Options{
		SourceType:      STATIC,
		RefreshInterval: 0,
		FileSys:         os.DirFS("."),
	}

	// copy options
	for _, o := range opts {
		o(&options)
	}

	switch sType := options.SourceType; sType {
	case KUBERNETES:
		panic("implement me")
	case DNS:
		panic("implement me")
	case FILE:
		r = &fileRefresher{
			version:         0,
			refreshInterval: options.RefreshInterval,
			fs:              options.FileSys,
			ctx:             ctx,
			snapshotCache:   cachev3.NewSnapshotCache(true, cachev3.IDHash{}, newXdsLogger()),
		}
	default:

		r = &staticRefresher{
			version:         0,
			refreshInterval: options.RefreshInterval,
			fs:              options.FileSys,
			ctx:             ctx,
			snapshotCache:   cachev3.NewSnapshotCache(true, cachev3.IDHash{}, newXdsLogger()),
		}
	}

	return
}
