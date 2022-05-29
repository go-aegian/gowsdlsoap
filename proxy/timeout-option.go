package proxy

import (
	"time"
)

// WithTimeout is an Option to set default HTTP dial context timeout.
func WithTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}
