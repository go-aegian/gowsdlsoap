package wsdl

import (
	"encoding/xml"
)

const wsdlNamespace = "http://schemas.xmlsoap.org/wsdl/"

// WSDL represents the global structure of file.
type WSDL struct {
	Xmlns           map[string]string `xml:"-"`
	Name            string            `xml:"name,attr"`
	TargetNamespace string            `xml:"targetNamespace,attr"`
	Imports         []*Import         `xml:"import"`
	Doc             string            `xml:"documentation"`
	Types           Type              `xml:"http://schemas.xmlsoap.org/wsdl/ types"`
	Messages        []*Message        `xml:"http://schemas.xmlsoap.org/wsdl/ message"`
	PortTypes       []*PortType       `xml:"http://schemas.xmlsoap.org/wsdl/ portType"`
	Binding         []*Binding        `xml:"http://schemas.xmlsoap.org/wsdl/ binding"`
	Service         []*Service        `xml:"http://schemas.xmlsoap.org/wsdl/ service"`
}

// UnmarshalXML implements interface xml.Unmarshaler for XSDSchema.
func (w *WSDL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	w.Xmlns = make(map[string]string)

	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			w.Xmlns[attr.Name.Local] = attr.Value
			continue
		}

		switch attr.Name.Local {
		case "name":
			w.Name = attr.Value

		case "targetNamespace":
			w.TargetNamespace = attr.Value
		}
	}

Loop:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch {
			case t.Name.Local == "import":
				x := new(Import)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				w.Imports = append(w.Imports, x)

			case t.Name.Local == "documentation":
				if err := d.DecodeElement(&w.Doc, &t); err != nil {
					return err
				}

			case t.Name.Space == wsdlNamespace:

				switch t.Name.Local {
				case "types":
					if err := d.DecodeElement(&w.Types, &t); err != nil {
						return err
					}
					for prefix, namespace := range w.Xmlns {
						for _, s := range w.Types.Schemas {
							if _, ok := s.Xmlns[prefix]; !ok {
								s.Xmlns[prefix] = namespace
							}
						}
					}

				case "message":
					x := new(Message)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}

					w.Messages = append(w.Messages, x)

				case "portType":
					x := new(PortType)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}

					w.PortTypes = append(w.PortTypes, x)

				case "binding":
					x := new(Binding)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}

					w.Binding = append(w.Binding, x)

				case "service":
					x := new(Service)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}

					w.Service = append(w.Service, x)

				default:
					err := d.Skip()
					if err != nil {
						return err
					}

					continue Loop
				}

			default:
				err := d.Skip()
				if err != nil {
					return err
				}
				continue Loop
			}

		case xml.EndElement:
			break Loop
		}
	}

	return nil
}
