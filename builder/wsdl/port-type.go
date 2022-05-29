package wsdl

// PortType defines the service, operations that can be performed and the messages involved.
// A port type can be compared to a function library, module or class.
type PortType struct {
	Name       string       `xml:"name,attr"`
	Doc        string       `xml:"documentation"`
	Operations []*Operation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}
