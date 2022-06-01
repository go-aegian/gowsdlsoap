package xsd

// Attribute represent an element attribute. Simple elements cannot have
// attributes. If an element has attributes, it is considered to be of a
// complex type. But the attribute itself is always declared as a simple type.
type Attribute struct {
	Doc        string      `xml:"annotation>documentation"`
	Name       string      `xml:"name,attr"`
	Ref        string      `xml:"ref,attr"`
	Type       string      `xml:"type,attr"`
	Use        string      `xml:"use,attr"`
	Fixed      string      `xml:"fixed,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
	Abstract   bool        `xml:"abstract,attr"`
}
