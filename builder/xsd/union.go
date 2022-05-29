package xsd

// Union represents a union element
type Union struct {
	SimpleType  []*SimpleType `xml:"simpleType,omitempty"`
	MemberTypes string        `xml:"memberTypes,attr"`
}
