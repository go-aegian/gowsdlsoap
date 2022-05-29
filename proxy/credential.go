package proxy

// Credential are used when setting WithBasicAuth or WithNTLM,
// the latter needs the domain to be set, whereas for basic auth is not needed.
type Credential struct {
	Username string
	Password string
}

type DomainCredential struct {
	Domain string
	Credential
}
