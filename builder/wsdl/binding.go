package wsdl

// Binding defines only a SOAP binding and its operations
type Binding struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	Doc         string       `xml:"documentation"`
	SOAPBinding SOAPBinding  `xml:"http://schemas.xmlsoap.org/wsdl/soap/ binding"`
	Operations  []*Operation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}
