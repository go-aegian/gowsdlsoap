package soap

import "encoding/xml"

type HeaderResponse struct {
	XMLName xml.Name `xml:"soap:Header"`
	Headers []interface{}
}
