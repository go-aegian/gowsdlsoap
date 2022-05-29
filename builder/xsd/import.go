package xsd

import "encoding/xml"

// Import represents XSD imports within the main schema.
type Import struct {
	XMLName        xml.Name `xml:"import"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Namespace      string   `xml:"namespace,attr"`
}
