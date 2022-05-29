package soap

import "encoding/xml"

type Header struct {
	XMLName xml.Name `xml:"soap:Header"`
	Headers []interface{}
}
