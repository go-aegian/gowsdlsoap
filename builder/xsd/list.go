package xsd

// List represents a element list
type List struct {
	Doc        string      `xml:"annotation>documentation"`
	ItemType   string      `xml:"itemType,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
}
