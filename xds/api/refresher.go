package api

import "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

type Refresher interface {
	GetCache() cache.Cache
	Start() error
}
