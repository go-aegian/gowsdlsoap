package proxy

import (
	"crypto/tls"
)

// WithTLS is an Option to set tls config
// This option cannot be used with WithHTTPClient
func WithTLS(tls *tls.Config) Option {
	return func(o *Options) {
		o.TlsConfig = tls
	}
}
