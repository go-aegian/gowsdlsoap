package proxy

import (
	"time"
)

// WithTLSHandshakeTimeout is an Option to set default tls handshake timeout.
//
// This option cannot be used with WithHTTPClient.
func WithTLSHandshakeTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.TlsHandshakeTimeout = t
	}
}
