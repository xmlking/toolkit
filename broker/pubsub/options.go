package broker

import (
	"context"

	"google.golang.org/api/option"
)

// https://github.com/cloudevents/sdk-go/blob/master/protocol/pubsub/v2/options.go

type Option func(*Options)

type Options struct {
	Name          string
	ClientOptions []option.ClientOption
	ProjectID     string

	// Handler executed when error happens in broker message
	// processing
	ErrorHandler Handler

	Context context.Context
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

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

// Context  appContext to trigger terminate signal
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
