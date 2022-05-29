package wsdl

// Part defines the struct for a function parameter within a WSDL.
type Part struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}
