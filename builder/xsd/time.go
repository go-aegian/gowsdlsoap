package xsd

import (
	"encoding/xml"
	"strings"
	"time"
)

// Time is a type for representing xsd:time
type Time struct {
	InnerTime time.Time
	hasTz     bool
}

// NewTime creates an object representing xsd:time in Golang
func NewTime(hour int, min int, sec int, nanoseconds int, loc *time.Location) Time {
	realLoc := loc
	if realLoc == nil {
		realLoc = time.Local
	}

	return Time{InnerTime: time.Date(1951, 10, 22, hour, min, sec, nanoseconds, realLoc), hasTz: loc != nil}
}

// MarshalXML implements xml.MarshalAttr on XSDTime
func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xdtString := t.string()
	if xdtString != "" {
		return e.EncodeElement(xdtString, start)
	}

	return nil
}

// MarshalXMLAttr implements xml.MarshalAttr on XSDTime
func (t Time) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	xdString := t.string()
	attr := xml.Attr{}

	if xdString != "" {
		attr.Name = name
		attr.Value = xdString
	}

	return attr, nil
}

// returns string representation and skips "zero" time values
func (t Time) string() string {
	if t.InnerTime.IsZero() {
		return ""
	}

	dateTimeLayout := time.RFC3339Nano

	if t.InnerTime.Nanosecond() == 0 {
		dateTimeLayout = time.RFC3339
	}

	// split off date portion
	dateAndTime := strings.SplitN(t.InnerTime.Format(dateTimeLayout), "T", 2)
	timeString := dateAndTime[1]
	if !t.hasTz {
		runes := strings.SplitN(timeString, "Z", 2)
		runes = strings.SplitN(runes[0], "+", 2)
		runes = strings.SplitN(runes[0], "-", 2)
		timeString = runes[0]
	}

	return timeString

}

// UnmarshalXML implements xml.Unmarshaler on XSDTime to use dateTimeLayout
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string

	err := d.DecodeElement(&content, &start)
	if err != nil {
		return err
	}

	return t.fromString(content)
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr on XSDTime to use dateTimeLayout
func (t *Time) UnmarshalXMLAttr(attr xml.Attr) error {
	return t.fromString(attr.Value)
}

func (t *Time) fromString(content string) error {
	var newTime time.Time
	var err error

	if content == "" {
		t.InnerTime = newTime
		return nil
	}

	t.hasTz = strings.Contains(content, "Z") || strings.Contains(content, "+") || strings.Contains(content, "-")

	if !t.hasTz {
		content += "Z"
	}

	t.InnerTime, err = time.Parse(timeLayout, content)

	return err
}

// Hour returns hour of the xsd:time
func (t Time) Hour() int {
	return t.InnerTime.Hour()
}

// Minute returns minutes of the xsd:time
func (t Time) Minute() int {
	return t.InnerTime.Minute()
}

// Second returns seconds of the xsd:time
func (t Time) Second() int {
	return t.InnerTime.Second()
}

// Nanosecond returns nanosecond of the xsd:time
func (t Time) Nanosecond() int {
	return t.InnerTime.Nanosecond()
}

// Location returns the TZ information of the xsd:time
func (t Time) Location() *time.Location {
	if t.hasTz {
		return t.InnerTime.Location()
	}
	return nil
}
