package wsdl

// Input represents input message.
type Input struct {
	Name       string        `xml:"name,attr"`
	Message    string        `xml:"message,attr"`
	Doc        string        `xml:"documentation"`
	SOAPBody   SOAPBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SOAPHeader []*SOAPHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
}
