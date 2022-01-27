package xds

import (
	"io/fs"
	"time"
)

type Option func(*Options)

type Options struct {
	//SourceType      string
	RefreshInterval time.Duration
	NodeID          string
	Namespace       string
	Hostnames       []string
	// FileSystem to load config files from. default: os.DirFS(".")
	FileSys fs.FS
}

// SourceType Type of the endpoints source
//func SourceType(t string) Option {
//	return func(o *Options) {
//		o.SourceType = t
//	}
//}

// WithRefreshInterval specifies the interval to poll Source for endpoints updates. default = 0, means: never refresh
func WithRefreshInterval(interval time.Duration) Option {
	return func(o *Options) {
		if interval <= 0 {
			o.RefreshInterval = 0
		} else {
			o.RefreshInterval = interval
		}
	}
}

// WithNodeID set NodeID in UUID format which should math to xDS Bootstrap file
func WithNodeID(nodeID string) Option {
	return func(o *Options) {
		o.NodeID = nodeID
	}
}

// *** for file source *** //

// WithFS enables use custom FileSystem to load config files. e.g., embed.FS
// default: os.DirFS(".")
func WithFS(fs fs.FS) Option {
	return func(o *Options) {
		o.FileSys = fs
	}
}

// *** for kubernetes source *** //

// WithNamespace : kubernetes namespace to monitor for endpoints
func WithNamespace(namespace string) Option {
	return func(o *Options) {
		o.Namespace = namespace
	}
}

// WithHostnames : kubernetes namespace to monitor for endpoints
func WithHostnames(hostnames []string) Option {
	return func(o *Options) {
		o.Hostnames = hostnames
	}
}
