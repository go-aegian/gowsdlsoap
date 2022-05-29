package xsd

import "encoding/xml"

// ComplexType represents a Schema complex type.
type ComplexType struct {
	XMLName        xml.Name       `xml:"complexType"`
	Abstract       bool           `xml:"abstract,attr"`
	Name           string         `xml:"name,attr"`
	Mixed          bool           `xml:"mixed,attr"`
	Sequence       []*Element     `xml:"sequence>element"`
	Choice         []*Element     `xml:"choice>element"`
	SequenceChoice []*Element     `xml:"sequence>choice>element"`
	All            []*Element     `xml:"all>element"`
	ComplexContent ComplexContent `xml:"complexContent"`
	SimpleContent  SimpleContent  `xml:"simpleContent"`
	Attributes     []*Attribute   `xml:"attribute"`
	Any            []*Any         `xml:"sequence>any"`
}
