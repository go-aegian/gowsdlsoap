package proxy

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-aegian/gosoap/builder/soap"
)

type mtomDecoder struct {
	reader *multipart.Reader
}

func getMtomHeader(contentType string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			return "", fmt.Errorf("Invalid multipart boundary: %s", boundary)
		}

		cType, ok := params["type"]
		if !ok || cType != "application/xop+xml" {
			// Process as normal xml (Don't resolve XOP parts)
			return "", nil
		}

		startInfo, ok := params["start-info"]
		if !ok || startInfo != "application/soap+xml" {
			return "", fmt.Errorf(`Expected param start-info="application/soap+xml", got %s`, startInfo)
		}
		return boundary, nil
	}

	return "", nil
}

func newMtomDecoder(r io.Reader, boundary string) *mtomDecoder {
	return &mtomDecoder{
		reader: multipart.NewReader(r, boundary),
	}
}

func (d *mtomDecoder) Decode(v interface{}) error {
	fields := make([]reflect.Value, 0)
	getBinaryFields(v, &fields)

	packages := make(map[string]*Binary, 0)
	for {
		p, err := d.reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		contentType := p.Header.Get(soap.ContentTypeHeader)
		if contentType == "application/xop+xml" {
			err := xml.NewDecoder(p).Decode(v)
			if err != nil {
				return err
			}
		} else {
			contentID := p.Header.Get(soap.ContentIdHeader)
			if contentID == "" {
				return errors.New("Invalid multipart content ID")
			}

			content, err := ioutil.ReadAll(p)
			if err != nil {
				return err
			}

			contentID = strings.Trim(contentID, "<>")
			packages[contentID] = &Binary{
				content:     &content,
				contentType: contentType,
			}
		}
	}

	// Set binary fields with correct content
	for _, f := range fields {
		b := f.Interface().(*Binary)
		b.content = packages[b.packageID].content
		b.contentType = packages[b.packageID].contentType
	}

	return nil
}
