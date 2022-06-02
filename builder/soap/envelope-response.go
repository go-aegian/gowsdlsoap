package soap

import (
	"encoding/xml"
	"strings"
)

type EnvelopeResponse struct {
	XMLName     xml.Name   `xml:"Envelope"`
	Attr        []xml.Attr `xml:",any,attr,omitempty"`
	Header      *HeaderResponse
	Body        BodyResponse
	Attachments []MIMEMultipartAttachment `xml:"attachments,omitempty"`
}

func NewEnvelopeResponse(ns map[string]string) *EnvelopeResponse {
	env := &EnvelopeResponse{}

	env.addXmlns("xmlns:soap", XmlNsSoapEnv)

	// env.setXmlns(ns)

	return env
}

func (e *EnvelopeResponse) setXmlns(ns map[string]string) {
	for alias, value := range ns {
		e.addXmlns(alias, value)
	}
}

func (e *EnvelopeResponse) addXmlns(alias, value string) {
	if !strings.HasPrefix(alias, "xmlns:") {
		alias = "xmlns:" + alias
	}
	e.Attr = append(e.Attr, xml.Attr{Name: xml.Name{Local: alias}, Value: value})
}
