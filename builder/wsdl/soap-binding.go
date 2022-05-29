package wsdl

// SOAPBinding represents a SOAP binding to the web service.
type SOAPBinding struct {
	Style     string `xml:"style,attr"`
	Transport string `xml:"transport,attr"`
}
