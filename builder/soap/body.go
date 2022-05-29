package soap

import "encoding/xml"

type Body struct {
	XMLName xml.Name    `xml:"soap:Body"`
	Content interface{} `xml:",omitempty"`
	Fault   *Fault      `xml:",omitempty"`
	faulted bool
}

func (b *Body) ErrorFromFault() error {
	if b.faulted {
		return b.Fault
	}
	b.Fault = nil
	return nil
}
