package xsd

import (
	"encoding/xml"
	"strings"
	"time"
)

// Date is a type for representing xsd:date in Golang
type Date struct {
	innerDate time.Time
	hasTz     bool
}

// NewDate creates an object represent xsd:datetime object in Golang
func NewDate(date time.Time, hasTz bool) Date {
	return Date{innerDate: date, hasTz: hasTz}
}

// StripTz removes the TZ information from the date
func (d *Date) StripTz() {
	d.hasTz = false
}

// Time converts the date to time.Time.
func (d *Date) Time() time.Time {
	if d.hasTz {
		return d.innerDate
	}

	return time.Date(d.innerDate.Year(), d.innerDate.Month(), d.innerDate.Day(), 0, 0, 0, 0, time.Local)
}

// MarshalXML implementation on XSDDate
func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xdtString := d.string()
	if xdtString != "" {
		return e.EncodeElement(xdtString, start)
	}

	return nil
}

// MarshalXMLAttr implementation on XSDDate
func (d Date) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	xdString := d.string()
	attr := xml.Attr{}

	if xdString != "" {
		attr.Name = name
		attr.Value = xdString
	}

	return attr, nil
}

// returns string representation and skips "zero" time values
func (d Date) string() string {
	if d.innerDate.IsZero() {
		return ""
	}

	dateString := d.innerDate.Format(dateLayout) // serialize with TZ

	if !d.hasTz {
		if strings.Contains(dateString, "Z") {
			// UTC Tz
			runes := strings.SplitN(dateString, "Z", 2)
			dateString = runes[0]
		} else {
			// [+-]00:00 Tz, remove last 6 chars
			if len(dateString) > 5 { // this should always be true
				start := len(dateString) - 6 // locate at "-"
				dateString = dateString[0:start]
			}
		}
	}

	return dateString
}

// UnmarshalXML implements xml.Unmarshaler on Date to use dateLayout
func (d *Date) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string

	err := decoder.DecodeElement(&content, &start)
	if err != nil {
		return err
	}

	d.innerDate, d.hasTz, err = fromString(content, dateLayout)

	return err
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr on XSDDate to use dateLayout
func (d *Date) UnmarshalXMLAttr(attr xml.Attr) error {
	var err error

	d.innerDate, d.hasTz, err = fromString(attr.Value, dateLayout)

	return err
}
