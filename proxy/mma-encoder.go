package proxy

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/go-aegian/gowsdlsoap/builder/soap"
)

const mmaContentType string = `multipart/related; start="<soap-request@gowsdlsoap.proxy>"; type="text/xml"; boundary="%s"`

type mmaEncoder struct {
	writer      *multipart.Writer
	attachments []soap.MIMEMultipartAttachment
}

func newMmaEncoder(w io.Writer, attachments []soap.MIMEMultipartAttachment) *mmaEncoder {
	return &mmaEncoder{
		writer:      multipart.NewWriter(w),
		attachments: attachments,
	}
}

func (e *mmaEncoder) Encode(v interface{}) error {
	var err error
	var soapPartWriter io.Writer

	// 1. write SOAP envelope part
	headers := make(textproto.MIMEHeader)
	headers.Set(soap.ContentTypeHeader, `text/xml;charset=UTF-8`)
	headers.Set(soap.ContentTransferEncodingHeader, "8bit")
	headers.Set(soap.ContentIdHeader, "<soap-request@gowsdlsoap.proxy>")
	if soapPartWriter, err = e.writer.CreatePart(headers); err != nil {
		return err
	}
	xmlEncoder := xml.NewEncoder(soapPartWriter)
	if err := xmlEncoder.Encode(v); err != nil {
		return err
	}

	// 2. write attachments parts
	for _, attachment := range e.attachments {
		attHeader := make(textproto.MIMEHeader)
		attHeader.Set(soap.ContentTypeHeader, fmt.Sprintf("application/octet-stream; name=%s", attachment.Name))
		attHeader.Set(soap.ContentTransferEncodingHeader, "binary")
		attHeader.Set(soap.ContentIdHeader, fmt.Sprintf("<%s>", attachment.Name))
		attHeader.Set("Content-Disposition",
			fmt.Sprintf("attachment; name=\"%s\"; filename=\"%s\"", attachment.Name, attachment.Name))

		attachmentPartWriter, err := e.writer.CreatePart(attHeader)
		if err != nil {
			return err
		}

		_, err = io.Copy(attachmentPartWriter, bytes.NewReader(attachment.Data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *mmaEncoder) Flush() error {
	return e.writer.Close()
}

func (e *mmaEncoder) Boundary() string {
	return e.writer.Boundary()
}

func getMmaHeader(contentType string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			return "", fmt.Errorf("invalid multipart boundary: %s", boundary)
		}

		startInfo, ok := params["start"]
		if !ok || startInfo != "<soap-request@gowsdlsoap.proxy>" {
			return "", fmt.Errorf(`expected param start="<soap-request@gowsdlsoap.proxy>", got %s`, startInfo)
		}
		return boundary, nil
	}

	return "", nil
}
