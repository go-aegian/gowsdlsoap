package proxy

import (
	"crypto/tls"
	"time"

	"github.com/vadimi/go-http-ntlm/v2"
)

type Options struct {
	Client              HTTPClient
	Transport           *httpntlm.NtlmTransport
	TlsConfig           *tls.Config
	HttpHeaders         map[string]string
	Mtom                bool
	Mma                 bool
	Timeout             time.Duration
	ConnectionTimeout   time.Duration
	TlsHandshakeTimeout time.Duration
	BasicAuth           *Credential
	NtlmAuth            *DomainCredential
}

// Option allows to customize the default http client or
// even provide a custom http client.
type Option func(*Options)

var DefaultOptions = Options{
	Timeout:             30 * time.Second,
	ConnectionTimeout:   90 * time.Second,
	TlsHandshakeTimeout: 15 * time.Second,
}
