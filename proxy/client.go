package proxy

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/go-aegian/gosoap/builder/soap"
)

// HTTPClient is a client which can make HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client - soap client
type Client struct {
	url         string
	opts        *Options
	headers     []interface{}
	attachments []soap.MIMEMultipartAttachment
}

// NewClient creates new SOAP client instance
func NewClient(url string, opt ...Option) *Client {
	opts := DefaultOptions
	for _, o := range opt {
		o(&opts)
	}

	return &Client{url: url, opts: &opts}
}

// AddHeader adds envelope header
// For correct behavior, every header must contain a `XMLName` field.  Refer to #121 for details
func (s *Client) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

// AddMIMEMultipartAttachment adds an attachment to the client that will be sent only if the
// WithMIMEMultipartAttachments option is used
func (s *Client) AddMIMEMultipartAttachment(attachment soap.MIMEMultipartAttachment) {
	s.attachments = append(s.attachments, attachment)
}

// SetHeaders sets envelope headers, overwriting any existing headers.
// For correct behavior, every header must contain a `XMLName` field.  Refer to #121 for details
func (s *Client) SetHeaders(headers ...interface{}) {
	s.headers = headers
}

// CallContext performs HTTP POST request with a context
func (s *Client) CallContext(ctx context.Context, soapAction string, request, response interface{}) error {
	return s.call(ctx, soapAction, request, response, nil, nil)
}

// Call performs HTTP POST request.
// Note that if the server returns a status code >= 400, a HTTPError will be returned
func (s *Client) Call(soapAction string, request, response interface{}) error {
	return s.call(context.Background(), soapAction, request, response, nil, nil)
}

// CallContextWithAttachmentsAndFaultDetail performs HTTP POST request.
// Note that if SOAP fault is returned, it will be stored in the error.
// On top the attachments array will be filled with attachments returned from the SOAP request.
func (s *Client) CallContextWithAttachmentsAndFaultDetail(ctx context.Context, soapAction string, request, response interface{}, faultDetail soap.FaultError, attachments *[]soap.MIMEMultipartAttachment) error {
	return s.call(ctx, soapAction, request, response, faultDetail, attachments)
}

// CallContextWithFaultDetail performs HTTP POST request.
// Note that if SOAP fault is returned, it will be stored in the error.
func (s *Client) CallContextWithFaultDetail(ctx context.Context, soapAction string, request, response interface{}, faultDetail soap.FaultError) error {
	return s.call(ctx, soapAction, request, response, faultDetail, nil)
}

// CallWithFaultDetail performs HTTP POST request.
// Note that if SOAP fault is returned, it will be stored in the error.
// the passed in fault detail is expected to implement FaultError interface,
// which allows to condense the detail into a short error message.
func (s *Client) CallWithFaultDetail(soapAction string, request, response interface{}, faultDetail soap.FaultError) error {
	return s.call(context.Background(), soapAction, request, response, faultDetail, nil)
}

func (s *Client) call(ctx context.Context, soapAction string, request, response interface{}, faultDetail soap.FaultError,
	retAttachments *[]soap.MIMEMultipartAttachment) error {

	soapRequest := soap.NewEnvelope()
	defer LogXml("Request", soapRequest)

	if s.headers != nil && len(s.headers) > 0 {
		soapRequest.Header = &soap.Header{Headers: s.headers}
	}

	soapRequest.Body.Content = request

	buffer := new(bytes.Buffer)

	var encoder soap.Encoder
	if s.opts.Mtom && s.opts.Mma {
		return fmt.Errorf("cannot use MTOM (XOP) and MMA (MIME Multipart Attachments) option at the same time")
	}

	if s.opts.Mtom {
		encoder = newMtomEncoder(buffer)
	} else if s.opts.Mma {
		encoder = newMmaEncoder(buffer, s.attachments)
	} else {
		encoder = xml.NewEncoder(buffer)
	}

	if err := encoder.Encode(soapRequest); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, s.url, buffer)
	if err != nil {
		return err
	}

	if s.opts.BasicAuth != nil {
		httpRequest.SetBasicAuth(s.opts.BasicAuth.Username, s.opts.BasicAuth.Password)
	}

	httpRequest = httpRequest.WithContext(ctx)

	if s.opts.Mtom {
		httpRequest.Header.Add(soap.ContentTypeHeader, fmt.Sprintf(soap.MtomContentType, encoder.(*mtomEncoder).Boundary()))
	} else if s.opts.Mma {
		httpRequest.Header.Add(soap.ContentTypeHeader, fmt.Sprintf(mmaContentType, encoder.(*mmaEncoder).Boundary()))
	} else {
		httpRequest.Header.Add(soap.ContentTypeHeader, "text/xml; charset=\"utf-8\"")
	}

	httpRequest.Header.Add("SOAPAction", soapAction)
	httpRequest.Header.Set("User-Agent", "gosoap/1.0")

	if s.opts.HttpHeaders != nil {
		for k, v := range s.opts.HttpHeaders {
			httpRequest.Header.Set(k, v)
		}
	}

	httpRequest.Close = true

	client := s.opts.Client
	if client == nil {
		if s.opts.Transport != nil {
			s.opts.Transport.RoundTripper.(*http.Transport).TLSHandshakeTimeout = s.opts.TlsHandshakeTimeout
			if s.opts.TlsConfig != nil {
				s.opts.Transport.RoundTripper.(*http.Transport).TLSClientConfig = s.opts.TlsConfig
			}
			client = &http.Client{Timeout: s.opts.ConnectionTimeout, Transport: s.opts.Transport}
		} else {
			tr := &http.Transport{
				TLSClientConfig: s.opts.TlsConfig,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					d := net.Dialer{Timeout: s.opts.Timeout}
					return d.DialContext(ctx, network, addr)
				},
				TLSHandshakeTimeout: s.opts.TlsHandshakeTimeout,
			}
			client = &http.Client{Timeout: s.opts.ConnectionTimeout, Transport: tr}
		}
	}

	res, err := client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(res.Body)
		return &soap.HTTPError{
			StatusCode:   res.StatusCode,
			ResponseBody: body,
		}
	}

	// xml Decoder (used with and without MTOM) cannot handle namespace prefixes (yet),
	// so use a namespace-less response envelope
	soapResponse := soap.NewEnvelopeResponse()
	soapResponse.Body = soap.BodyResponse{
		Content: response,
		Fault: &soap.Fault{
			Detail: faultDetail,
		},
	}
	defer LogXml("Response", soapResponse)

	mtomBoundary, err := getMtomHeader(res.Header.Get(soap.ContentTypeHeader))
	if err != nil {
		return err
	}

	var mmaBoundary string
	if s.opts.Mma {
		mmaBoundary, err = getMmaHeader(res.Header.Get(soap.ContentTypeHeader))
		if err != nil {
			return err
		}
	}

	var dec soap.Decoder
	if mtomBoundary != "" {
		dec = newMtomDecoder(res.Body, mtomBoundary)
	} else if mmaBoundary != "" {
		dec = newMmaDecoder(res.Body, mmaBoundary)
	} else {
		dec = xml.NewDecoder(res.Body)
	}

	if err := dec.Decode(soapResponse); err != nil {
		return err
	}

	if soapResponse.Attachments != nil {
		*retAttachments = soapResponse.Attachments
	}
	return soapResponse.Body.ErrorFromFault()
}

func LogXml(logType string, message interface{}) {
	marshalledRequest, err := xml.MarshalIndent(message, "", "\t")
	if err != nil {
		log.Fatalf("\nerror parsing as xml: %s %v %v\n", logType, message, err)
	}

	fmt.Printf("\n%s:\n%s\n\n", logType, string(marshalledRequest))
}
