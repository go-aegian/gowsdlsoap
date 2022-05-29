package proxy

import (
	"encoding/xml"
)

const (
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName        xml.Name          `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse      string            `xml:"xmlns:wsse,attr"`
	MustUnderstand string            `xml:"mustUnderstand,attr,omitempty"`
	Token          *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name     `xml:"wsse:UsernameToken"`
	XmlNSWsu  string       `xml:"xmlns:wsu,attr"`
	XmlNSWsse string       `xml:"xmlns:wsse,attr"`
	Id        string       `xml:"wsu:Id,attr,omitempty"`
	Username  *WSSUsername `xml:",omitempty"`
	Password  *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	Data      string   `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`
	Data      string   `xml:",chardata"`
}

// NewWSSSecurityHeader creates WSSSecurityHeader instance
func NewWSSSecurityHeader(credential *Credential, tokenID, mustUnderstand string) *WSSSecurityHeader {
	return &WSSSecurityHeader{
		XmlNSWsse:      WssNsWSSE,
		MustUnderstand: mustUnderstand,
		Token: &WSSUsernameToken{
			XmlNSWsu:  WssNsWSU,
			XmlNSWsse: WssNsWSSE,
			Id:        tokenID,
			Username: &WSSUsername{
				XmlNSWsse: WssNsWSSE,
				Data:      credential.Username,
			},
			Password: &WSSPassword{
				XmlNSWsse: WssNsWSSE,
				XmlNSType: WssNsType,
				Data:      credential.Password,
			},
		},
	}
}
