package wsdl

// Fault represents a fault message.
type Fault struct {
	Name      string    `xml:"name,attr"`
	Message   string    `xml:"message,attr"`
	Doc       string    `xml:"documentation"`
	SOAPFault SOAPFault `xml:"http://schemas.xmlsoap.org/wsdl/soap/ fault"`
}
