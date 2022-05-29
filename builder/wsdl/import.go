package wsdl

// Import is the struct used for deserializing WSDL imports.
type Import struct {
	Namespace string `xml:"namespace,attr"`
	Location  string `xml:"location,attr"`
}
