package xsd

import "encoding/xml"

// Element represents a Schema element.
type Element struct {
	XMLName     xml.Name     `xml:"element"`
	Name        string       `xml:"name,attr"`
	Doc         string       `xml:"annotation>documentation"`
	Nillable    bool         `xml:"nillable,attr"`
	Type        string       `xml:"type,attr"`
	Ref         string       `xml:"ref,attr"`
	MinOccurs   string       `xml:"minOccurs,attr"`
	MaxOccurs   string       `xml:"maxOccurs,attr"`
	ComplexType *ComplexType `xml:"complexType"` // local
	SimpleType  *SimpleType  `xml:"simpleType"`
	Groups      []*Group     `xml:"group"`
}
