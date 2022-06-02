package soap

import "encoding/xml"

type Envelope struct {
	XMLName  xml.Name `xml:"soap:Envelope"`
	XMLNS    string   `xml:"xmlns:soap,attr"`
	XMLNSXsd string   `xml:"xmlns:xsd,attr,omitempty"`
	XMLNSXsi string   `xml:"xmlns:xsi,attr,omitempty"`
	Header   *Header
	Body     Body
}

func NewEnvelope() *Envelope {
	return &Envelope{XMLNS: XmlNsSoapEnv, XMLNSXsd: XmlNsSoapXsd, XMLNSXsi: XmlNsSoapXsi}
}
