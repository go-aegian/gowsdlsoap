package xsd

import "encoding/xml"

const xmlschema11 = "http://www.w3.org/2001/XMLSchema"

// Schema represents an entire Schema structure.
type Schema struct {
	XMLName            xml.Name          `xml:"schema"`
	Xmlns              map[string]string `xml:"-"`
	Tns                string            `xml:"xmlns tns,attr"`
	Xs                 string            `xml:"xmlns xs,attr"`
	Version            string            `xml:"version,attr"`
	TargetNamespace    string            `xml:"targetNamespace,attr"`
	ElementFormDefault string            `xml:"elementFormDefault,attr"`
	Includes           []*Include        `xml:"include"`
	Imports            []*Import         `xml:"import"`
	Elements           []*Element        `xml:"element"`
	Attributes         []*Attribute      `xml:"attribute"`
	ComplexTypes       []*ComplexType    `xml:"complexType"`
	SimpleType         []*SimpleType     `xml:"simpleType"`
}

// UnmarshalXML implements interface xml.Unmarshaler for XSDSchema.
func (s *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s.Xmlns = make(map[string]string)
	s.XMLName = start.Name

	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			s.Xmlns[attr.Name.Local] = attr.Value
			continue
		}

		switch attr.Name.Local {
		case "version":
			s.Version = attr.Value
		case "targetNamespace":
			s.TargetNamespace = attr.Value
		case "elementFormDefault":
			s.ElementFormDefault = attr.Value
		}
	}

Loop:
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Space != xmlschema11 {
				err := d.Skip()
				if err != nil {
					return err
				}
				continue Loop
			}

			switch t.Name.Local {
			case "include":
				x := new(Include)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.Includes = append(s.Includes, x)

			case "import":
				x := new(Import)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.Imports = append(s.Imports, x)

			case "element":
				x := new(Element)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.Elements = append(s.Elements, x)

			case "attribute":
				x := new(Attribute)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.Attributes = append(s.Attributes, x)

			case "complexType":
				x := new(ComplexType)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.ComplexTypes = append(s.ComplexTypes, x)

			case "simpleType":
				x := new(SimpleType)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}

				s.SimpleType = append(s.SimpleType, x)

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
