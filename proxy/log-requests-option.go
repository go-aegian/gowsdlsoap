package proxy

// WithLogRequests sets logging of requests, by default log is disabled
func WithLogRequests(on bool) Option {
	return func(o *Options) {
		o.LogRequests = on
	}
}
