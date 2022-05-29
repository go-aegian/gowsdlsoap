package wsdl

// Port defines the properties for a SOAP port only.
type Port struct {
	Name        string      `xml:"name,attr"`
	Binding     string      `xml:"binding,attr"`
	Doc         string      `xml:"documentation"`
	SOAPAddress SOAPAddress `xml:"http://schemas.xmlsoap.org/wsdl/soap/ address"`
}
