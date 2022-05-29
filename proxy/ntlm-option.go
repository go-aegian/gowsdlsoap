package proxy

import (
	"crypto/tls"
	"net/http"

	"github.com/vadimi/go-http-ntlm/v2"
)

// WithNTLM configures for a given http.Client its transport for NTLM,
// it works with WithTLS which should be called prior to it to set the right tls.Config, if not it will
// default to a simple tls.Config with InsecureSkipVerify set to true to avoid certs issues
// This setting will be configured even in the custom httpClient replacing existing transport with this one.
func WithNTLM(credential *DomainCredential) Option {
	return func(o *Options) {
		if o.TlsConfig == nil {
			o.TlsConfig = &tls.Config{
				InsecureSkipVerify: true,
				ClientAuth:         tls.NoClientCert,
			}
		}

		o.NtlmAuth = credential

		o.Transport = &httpntlm.NtlmTransport{
			Domain:       credential.Domain,
			User:         credential.Username,
			Password:     credential.Password,
			RoundTripper: &http.Transport{TLSClientConfig: o.TlsConfig},
		}
	}
}
