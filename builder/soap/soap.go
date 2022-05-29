package soap

const (
	XmlNsSoapXsi                  = "http://www.w3.org/2001/XMLSchema-instance"
	XmlNsSoapXsd                  = "http://www.w3.org/2001/XMLSchema"
	XmlNsSoapEnv                  = "http://schemas.xmlsoap.org/soap/envelope/"
	MtomContentType               = `multipart/related; start-info="application/soap+xml"; type="application/xop+xml"; boundary="%s"`
	ContentTypeHeader             = "Content-Type"
	ContentTransferEncodingHeader = "Content-Transfer-Encoding"
	ContentIdHeader               = "Content-ID"
)

type Encoder interface {
	Encode(v interface{}) error
	Flush() error
}

type Decoder interface {
	Decode(v interface{}) error
}
