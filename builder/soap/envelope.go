package soap

import (
	"encoding/xml"
	"strings"
)

type Envelope struct {
	XMLName xml.Name   `xml:"soap:Envelope"`
	Attr    []xml.Attr `xml:",any,attr,omitempty"`
	Header  *Header
	Body    Body
}

func NewEnvelope(ns map[string]string) *Envelope {
	env := &Envelope{}
	env.addXmlns("xmlns:soap", XmlNsSoapEnv)

	env.setXmlns(ns)

	return env
}

func (e *Envelope) setXmlns(ns map[string]string) {
	for alias, value := range ns {
		e.addXmlns(alias, value)
	}
}

func (e *Envelope) addXmlns(alias, value string) {
	if !strings.HasPrefix(alias, "xmlns:") {
		alias = "xmlns:" + alias
	}

	// avoid duplicate aliases and urls namespaces
	for _, attribute := range e.Attr {
		if alias == attribute.Name.Local || value == attribute.Value {
			return
		}
	}

	e.Attr = append(e.Attr, xml.Attr{Name: xml.Name{Local: alias}, Value: value})
}
