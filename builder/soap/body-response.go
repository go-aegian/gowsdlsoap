package soap

import "encoding/xml"

type BodyResponse struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",omitempty"`
	Fault   *Fault      `xml:",omitempty"`
	faulted bool
}

// UnmarshalXML of the body xml
func (b *BodyResponse) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}
	if b.Fault == nil {
		b.Fault = &Fault{Detail: nil}
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			}
			if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Content = nil
				b.faulted = true

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (b *BodyResponse) ErrorFromFault() error {
	if b.faulted {
		return b.Fault
	}
	b.Fault = nil
	return nil
}
