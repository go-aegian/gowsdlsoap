package xsd

import "encoding/xml"

// Extension element extends an existing simpleType or complexType element.
type Extension struct {
	XMLName        xml.Name     `xml:"extension"`
	Base           string       `xml:"base,attr"`
	Attributes     []*Attribute `xml:"attribute"`
	Sequence       []*Element   `xml:"sequence>element"`
	Choice         []*Element   `xml:"choice>element"`
	SequenceChoice []*Element   `xml:"sequence>choice>element"`
}
