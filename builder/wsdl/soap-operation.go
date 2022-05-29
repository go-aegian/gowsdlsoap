package wsdl

// SOAPOperation represents a service operation in SOAP terms.
type SOAPOperation struct {
	SOAPAction string `xml:"soapAction,attr"`
	Style      string `xml:"style,attr"`
}
