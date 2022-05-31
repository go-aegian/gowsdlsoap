package proxy

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/go-aegian/gowsdlsoap/builder/soap"
)

type mmaDecoder struct {
	reader *multipart.Reader
}

func newMmaDecoder(r io.Reader, boundary string) *mmaDecoder {
	return &mmaDecoder{
		reader: multipart.NewReader(r, boundary),
	}
}

func (d *mmaDecoder) Decode(v interface{}) error {
	soapEnvResp := v.(*soap.EnvelopeResponse)
	attachments := make([]soap.MIMEMultipartAttachment, 0)
	for {
		p, err := d.reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if p.Header.Get(soap.ContentTypeHeader) == "text/xml;charset=UTF-8" {
			// decode SOAP part
			err = xml.NewDecoder(p).Decode(v)
			if err != nil {
				return err
			}
		} else {
			// decode attachment parts
			contentID := p.Header.Get(soap.ContentIdHeader)
			if contentID == "" {
				return errors.New("invalid multipart content id")
			}
			content, err := ioutil.ReadAll(p)
			if err != nil {
				return err
			}

			contentID = strings.Trim(contentID, "<>")
			attachments = append(attachments, soap.MIMEMultipartAttachment{
				Name: contentID,
				Data: content,
			})
		}
	}
	if len(attachments) > 0 {
		soapEnvResp.Attachments = attachments
	}

	return nil
}
