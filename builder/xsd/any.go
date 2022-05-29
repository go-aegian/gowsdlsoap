package xsd

import "encoding/xml"

// Any represents a Schema element.
type Any struct {
	XMLName         xml.Name `xml:"any"`
	Doc             string   `xml:"annotation>documentation"`
	MinOccurs       string   `xml:"minOccurs,attr"`
	MaxOccurs       string   `xml:"maxOccurs,attr"`
	Namespace       string   `xml:"namespace,attr"`
	ProcessContents string   `xml:"processContents,attr"`
}
