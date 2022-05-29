package xsd

import "encoding/xml"

// ComplexContent element defines extensions or restrictions on a complex
// type that contains mixed content or elements only.
type ComplexContent struct {
	XMLName   xml.Name  `xml:"complexContent"`
	Extension Extension `xml:"extension"`
}
