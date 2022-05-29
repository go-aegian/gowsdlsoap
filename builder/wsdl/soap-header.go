package wsdl

// SOAPHeader defines the header for a SOAP service.
type SOAPHeader struct {
	Message       string             `xml:"message,attr"`
	Part          string             `xml:"part,attr"`
	Use           string             `xml:"use,attr"`
	EncodingStyle string             `xml:"encodingStyle,attr"`
	Namespace     string             `xml:"namespace,attr"`
	HeadersFault  []*SOAPHeaderFault `xml:"headerfault"`
}
