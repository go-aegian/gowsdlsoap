package soap

import "encoding/xml"

type FaultError interface {
	// ErrorString should return a short version of the detail as a string,
	// which will be used in place of <faultstring> for the error message.
	// Set "HasData()" to always return false if <faultstring> error
	// message is preferred.
	ErrorString() string
	// HasData indicates whether the composite fault contains any data.
	HasData() bool
}

type Fault struct {
	XMLName xml.Name   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	Code    string     `xml:"faultcode,omitempty"`
	String  string     `xml:"faultstring,omitempty"`
	Actor   string     `xml:"faultactor,omitempty"`
	Detail  FaultError `xml:"detail,omitempty"`
}

func (f *Fault) Error() string {
	if f.Detail != nil && f.Detail.HasData() {
		return f.Detail.ErrorString()
	}
	return f.String
}
