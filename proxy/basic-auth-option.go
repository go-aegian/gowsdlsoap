package proxy

// WithBasicAuth is an Option to set BasicAuth
func WithBasicAuth(credential *Credential) Option {
	return func(o *Options) {
		o.BasicAuth = credential
	}
}
