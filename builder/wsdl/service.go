package wsdl

// Service defines the list of SOAP services associated with the WSDL.
type Service struct {
	Name  string  `xml:"name,attr"`
	Doc   string  `xml:"documentation"`
	Ports []*Port `xml:"http://schemas.xmlsoap.org/wsdl/ port"`
}
