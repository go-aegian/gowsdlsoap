package soap

import "encoding/xml"

type EnvelopeResponse struct {
	XMLName     xml.Name `xml:"Envelope"`
	XMLNS       string   `xml:"xmlns:soap,attr"`
	XMLNSXsd    string   `xml:"xmlns:xsd,attr,omitempty"`
	XMLNSXsi    string   `xml:"xmlns:xsi,attr,omitempty"`
	Header      *HeaderResponse
	Body        BodyResponse
	Attachments []MIMEMultipartAttachment `xml:"attachments,omitempty"`
}

func NewEnvelopeResponse() *EnvelopeResponse {
	return &EnvelopeResponse{XMLNS: XmlNsSoapEnv, XMLNSXsd: XmlNsSoapXsd, XMLNSXsi: XmlNsSoapXsi}
}
