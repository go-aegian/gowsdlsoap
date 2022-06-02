package proxy

// WithLogResponses sets logging of responses, by default log is disabled
func WithLogResponses(on bool) Option {
	return func(o *Options) {
		o.LogResponses = on
	}
}
