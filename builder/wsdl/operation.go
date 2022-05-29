package wsdl

// Operation represents the contract of an entire operation or function.
type Operation struct {
	Name          string        `xml:"name,attr"`
	Doc           string        `xml:"documentation"`
	Input         Input         `xml:"input"`
	Output        Output        `xml:"output"`
	Faults        []*Fault      `xml:"fault"`
	SOAPOperation SOAPOperation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
}
