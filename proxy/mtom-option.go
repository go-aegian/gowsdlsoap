package proxy

// WithMTOM is an Option to set Message Transmission Optimization Mechanism.
// MTOM encodes fields of type Binary using XOP.
func WithMTOM() Option {
	return func(o *Options) {
		o.Mtom = true
	}
}
