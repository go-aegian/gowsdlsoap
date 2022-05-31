package wsdl

import (
	"github.com/go-aegian/gowsdlsoap/builder/xsd"
)

// Type represents the entry point for deserializing XSD schemas used by the WSDL file.
type Type struct {
	Doc     string        `xml:"documentation"`
	Schemas []*xsd.Schema `xml:"schema"`
}
