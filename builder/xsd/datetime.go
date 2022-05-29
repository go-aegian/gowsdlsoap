package xsd

import (
	"encoding/xml"
	"strings"
	"time"
)

const (
	dateLayout = "2006-01-02Z07:00"
	timeLayout = "15:04:05.999999999Z07:00"
)

// DateTime represents xsd:datetime
type DateTime struct {
	innerTime time.Time
	hasTz     bool
}

// NewDateTime creates an object representing xsd:datetime
func NewDateTime(dt time.Time, hasTz bool) DateTime {
	return DateTime{innerTime: dt, hasTz: hasTz}
}

// StripTz removes TZ information from the datetime
func (dt *DateTime) StripTz() {
	dt.hasTz = false
}

// Time converts the time to time.Time.
func (dt *DateTime) Time() time.Time {
	if dt.hasTz {
		return dt.innerTime
	}

	return time.Date(
		dt.innerTime.Year(),
		dt.innerTime.Month(),
		dt.innerTime.Day(),
		dt.innerTime.Hour(),
		dt.innerTime.Minute(),
		dt.innerTime.Second(),
		dt.innerTime.Nanosecond(),
		time.Local,
	)
}

// MarshalXML implements xml.MarshalAttr on XsdDateTime
func (dt DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xdtString := dt.string()
	if xdtString != "" {
		return e.EncodeElement(xdtString, start)
	}

	return nil
}

// MarshalXMLAttr implements xml.MarshalAttr on XsdDateTime
func (dt DateTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	xdtString := dt.string()
	attr := xml.Attr{}
	if xdtString != "" {
		attr.Name = name
		attr.Value = xdtString
	}

	return attr, nil
}

// returns string representation and skips "zero" time values. It also checks if nanoseconds and TZ exist.
func (dt DateTime) string() string {
	if !dt.innerTime.IsZero() {
		dateTimeLayout := time.RFC3339Nano
		if dt.innerTime.Nanosecond() == 0 {
			dateTimeLayout = time.RFC3339
		}

		dtString := dt.innerTime.Format(dateTimeLayout)

		if !dt.hasTz {
			// split off time portion
			dateAndTime := strings.SplitN(dtString, "T", 2)
			runes := strings.SplitN(dateAndTime[1], "Z", 2)
			runes = strings.SplitN(runes[0], "+", 2)
			runes = strings.SplitN(runes[0], "-", 2)
			dtString = dateAndTime[0] + "T" + runes[0]
		}

		return dtString
	}

	return ""
}

// UnmarshalXML implements xml.Unmarshaler on XsdDateTime to use time.RFC3339Nano
func (dt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string

	err := d.DecodeElement(&content, &start)
	if err != nil {
		return err
	}

	dt.innerTime, dt.hasTz, err = fromString(content, time.RFC3339Nano)

	return err
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr on XsdDateTime to use time.RFC3339Nano
func (dt *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {
	var err error
	dt.innerTime, dt.hasTz, err = fromString(attr.Value, time.RFC3339Nano)
	return err
}

func fromString(content string, format string) (time.Time, bool, error) {
	var t time.Time
	if content == "" {
		return t, true, nil
	}

	hasTz := false
	if strings.Contains(content, "T") { // check if we have a time portion
		dateAndTime := strings.SplitN(content, "T", 2)
		if len(dateAndTime) > 1 {
			if strings.Contains(dateAndTime[1], "Z") ||
				strings.Contains(dateAndTime[1], "+") ||
				strings.Contains(dateAndTime[1], "-") {
				hasTz = true
			}
		}

		if !hasTz {
			content += "Z"
		}

		if content == "0001-01-01T00:00:00Z" {
			return t, true, nil
		}

	} else {
		// does not have a time portion thus check timezone
		if strings.Contains(content, "Z") ||
			strings.Contains(content, ":") {
			hasTz = true
		}

		if !hasTz {
			content += "Z"
		}
	}

	t, err := time.Parse(format, content)
	return t, hasTz, err
}
