package xsd

import "encoding/xml"

// SimpleContent element contains extensions or restrictions on a text-only
// complex type or on a simple type as content and contains no elements.
type SimpleContent struct {
	XMLName   xml.Name  `xml:"simpleContent"`
	Extension Extension `xml:"extension"`
}
