package xsd

// SimpleType element defines a simple type and specifies the constraints
// and information about the values of attributes or text-only elements.
type SimpleType struct {
	Name        string      `xml:"name,attr"`
	Doc         string      `xml:"annotation>documentation"`
	Restriction Restriction `xml:"restriction"`
	List        List        `xml:"list"`
	Union       Union       `xml:"union"`
	Final       string      `xml:"final"`
}
