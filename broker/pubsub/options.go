package broker

import (
	"context"
	"google.golang.org/api/option"
)

// https://github.com/cloudevents/sdk-go/blob/master/protocol/pubsub/v2/options.go

type Options struct {
	ClientOptions []option.ClientOption
	ProjectID     string

	// Handler executed when error happens in broker message
	// processing
	ErrorHandler Handler

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

// ClientOption is a broker Option which allows google pubsub client options to be
// set for the client
func ClientOption(c ...option.ClientOption) Option {
	return func(o *Options) {
		o.ClientOptions = c
	}
}

// ProjectID provides an option which sets the google project id
func ProjectID(id string) Option {
	return func(o *Options) {
		o.ProjectID = id
	}
}

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
