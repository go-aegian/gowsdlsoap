package proxy

import (
	"net/http"
)

// WithHTTPClient allows to use a custom HTTP client instead of the default.
//
// WithNTLM option, if provided, will replace the http client transport.
//
// WithTLSHandshakeTimeout, WithTLS and/or WithTimeout clientOption will be discarded.
func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.Client = c
		if o.Transport != nil {
			o.Client.(*http.Client).Transport = o.Transport
		}
	}
}
