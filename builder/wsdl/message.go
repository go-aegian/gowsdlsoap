package wsdl

// Message represents a function, which in turn has one or more parameters.
type Message struct {
	Name  string  `xml:"name,attr"`
	Doc   string  `xml:"documentation"`
	Parts []*Part `xml:"http://schemas.xmlsoap.org/wsdl/ part"`
}
