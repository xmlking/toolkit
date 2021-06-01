package broker

type PublishOptions struct {
	// publishes msg to the topic asynchronously if set to true.
	// Default false. i.e., publishes synchronously(blocking)
	Async bool
}

type PublishOption func(*PublishOptions)

func PublishAsync(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.Async = b
	}
}
