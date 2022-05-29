package proxy

import (
	"time"
)

// WithRequestTimeout is an Option to set default end-end connection timeout.
//
// This option cannot be used with WithHTTPClient
func WithRequestTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ConnectionTimeout = t
	}
}
