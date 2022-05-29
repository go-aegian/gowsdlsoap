package proxy

// WithHTTPHeaders is an Option to set global HTTP headers for all requests
func WithHTTPHeaders(headers map[string]string) Option {
	return func(o *Options) {
		o.HttpHeaders = headers
	}
}
