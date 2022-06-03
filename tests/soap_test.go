package tests

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-aegian/gowsdlsoap/builder/soap"
	"github.com/go-aegian/gowsdlsoap/builder/xsd"
	"github.com/go-aegian/gowsdlsoap/proxy"
	"github.com/stretchr/testify/assert"
)

type Ping struct {
	XMLName xml.Name     `xml:"http://example.com/service.xsd Ping"`
	Request *PingRequest `xml:"request,omitempty"`
}

type PingRequest struct {
	Message    string        `xml:"Message,omitempty"`
	Attachment *proxy.Binary `xml:"Attachment,omitempty"`
}

type PingResponse struct {
	XMLName    xml.Name    `xml:"http://example.com/service.xsd PingResponse"`
	PingResult *PingResult `xml:"PingResult,omitempty"`
}

type PingResult struct {
	XMLName    xml.Name      `xml:"PingResult"`
	Message    string        `xml:"Message,omitempty"`
	Attachment *proxy.Binary `xml:"Attachment,omitempty"`
}

type AttachmentRequest struct {
	XMLName   xml.Name `xml:"http://example.com/service.xsd attachmentRequest"`
	Name      string   `xml:"name,omitempty"`
	ContentID string   `xml:"contentID,omitempty"`
}

func TestClient_Call(t *testing.T) {
	pingRequest := &Ping{Request: &PingRequest{Message: "Ping"}}
	pingResponse := &PingResponse{PingResult: &PingResult{Message: "Pong"}}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := xml.NewDecoder(r.Body).Decode(pingRequest.Request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		soapEnv := soap.NewEnvelopeResponse()
		soapEnv.Body.Content = &pingResponse

		response, err := xml.Marshal(soapEnv)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		_, _ = w.Write(response)
	}))
	defer ts.Close()

	client := proxy.NewClient(ts.URL)

	err := client.Call("GetData", pingRequest, pingResponse)
	assert.NoError(t, err)

	assert.Equal(t, "Pong", pingResponse.PingResult.Message)
}

func TestClient_Send_Correct_Headers(t *testing.T) {
	tests := []struct {
		action          string
		reqHeaders      map[string]string
		expectedHeaders map[string]string
	}{
		{
			"GetTrade",
			map[string]string{},
			map[string]string{
				"User-Agent":           "gowsdlsoap/1.0",
				"SOAPAction":           "GetTrade",
				soap.ContentTypeHeader: "text/xml; charset=\"utf-8\"",
			},
		},
		{
			"SaveTrade",
			map[string]string{"User-Agent": "soap/0.1"},
			map[string]string{
				"User-Agent": "soap/0.1",
				"SOAPAction": "SaveTrade",
			},
		},
		{
			"SaveTrade",
			map[string]string{soap.ContentTypeHeader: "text/xml; charset=\"utf-16\""},
			map[string]string{soap.ContentTypeHeader: "text/xml; charset=\"utf-16\""},
		},
	}

	var gotHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
	}))
	defer ts.Close()

	for _, test := range tests {
		client := proxy.NewClient(ts.URL, proxy.WithHTTPHeaders(test.reqHeaders))
		req := struct{}{}
		reply := struct{}{}
		_ = client.Call(test.action, req, reply)

		for k, v := range test.expectedHeaders {
			h := gotHeaders.Get(k)
			if h != v {
				t.Errorf("got %s wanted %s", h, v)
			}
		}
	}
}

func TestClient_Attachments_WithAttachmentResponse(t *testing.T) {
	req := &AttachmentRequest{Name: "UploadMyFilePlease", ContentID: "First_Attachment"}

	firstAtt := soap.MIMEMultipartAttachment{Name: "First_Attachment", Data: []byte(`foobar`)}
	secondAtt := soap.MIMEMultipartAttachment{Name: "Second_Attachment", Data: []byte(`tl;tr`)}

	reply := &AttachmentRequest{}
	retAttachments := make([]soap.MIMEMultipartAttachment, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		bodyBuf, _ := ioutil.ReadAll(r.Body)
		_, _ = w.Write(bodyBuf)
	}))
	defer ts.Close()

	client := proxy.NewClient(ts.URL, proxy.WithMIMEMultipartAttachments())
	client.AddMIMEMultipartAttachment(firstAtt)
	client.AddMIMEMultipartAttachment(secondAtt)

	err := client.CallContextWithAttachmentsAndFaultDetail(context.TODO(), "''", req, reply, nil, &retAttachments)
	assert.NoError(t, err)

	assert.Equal(t, req.ContentID, reply.ContentID)
	assert.Len(t, retAttachments, 2)
	assert.Equal(t, retAttachments[0], firstAtt)
	assert.Equal(t, retAttachments[1], secondAtt)
}

func TestClient_MTOM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		bodyBuf, _ := ioutil.ReadAll(r.Body)
		_, _ = w.Write(bodyBuf)
	}))
	defer ts.Close()

	client := proxy.NewClient(ts.URL, proxy.WithMTOM())
	req := &PingRequest{Attachment: proxy.NewBinary([]byte("Attached data")).SetContentType("text/plain")}
	reply := &PingRequest{}
	if err := client.Call("GetData", req, reply); err != nil {
		t.Fatalf("couln't call service: %v", err)
	}

	if !bytes.Equal(reply.Attachment.Bytes(), req.Attachment.Bytes()) {
		t.Errorf("got %s wanted %s", reply.Attachment.Bytes(), req.Attachment.Bytes())
	}

	if reply.Attachment.ContentType() != req.Attachment.ContentType() {
		t.Errorf("got %s wanted %s", reply.Attachment.Bytes(), req.Attachment.ContentType())
	}
}

type SimpleNode struct {
	Detail string      `xml:"Detail,omitempty"`
	Num    float64     `xml:"Num,omitempty"`
	Nested *SimpleNode `xml:"Nested,omitempty"`
}

func (s SimpleNode) ErrorString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%.2f: %s", s.Num, s.Detail))
	if s.Nested != nil {
		sb.WriteString("\n" + s.Nested.ErrorString())
	}
	return sb.String()
}

func (s SimpleNode) HasData() bool {
	return true
}

type Wrapper struct {
	Item    interface{} `xml:"SimpleNode"`
	hasData bool
}

func (w *Wrapper) HasData() bool {
	return w.hasData
}

func (w *Wrapper) ErrorString() string {
	switch w.Item.(type) {
	case soap.FaultError:
		return w.Item.(soap.FaultError).ErrorString()
	}
	return "default error"
}

func Test_SimpleNode(t *testing.T) {
	input := `<SimpleNode>
				  <Name>SimpleNode</Name>
				  <Detail>detail message</Detail>
				  <Num>6.005</Num>
			  </SimpleNode>`
	decoder := xml.NewDecoder(strings.NewReader(input))

	simple := &SimpleNode{}
	err := decoder.Decode(simple)
	assert.NoError(t, err)
	assert.EqualValues(t, &SimpleNode{Detail: "detail message", Num: 6.005}, simple)
}

func Test_Client_FaultDefault(t *testing.T) {
	tests := []struct {
		name          string
		hasData       bool
		wantErrString string
		fault         interface{}
		emptyFault    interface{}
	}{
		{
			name:          "Empty-WithFault",
			wantErrString: "default error",
			hasData:       true,
		},
		{
			name:          "Empty-NoFaultDetail",
			wantErrString: "Custom error message.",
			hasData:       false,
		},
		{
			name:          "SimpleNode",
			wantErrString: "7.70: detail message",
			hasData:       true,
			fault: &SimpleNode{
				Detail: "detail message",
				Num:    7.7,
			},
			emptyFault: &SimpleNode{},
		},
		{
			name:          "ArrayOfNode",
			wantErrString: "default error",
			hasData:       true,
			fault: &[]SimpleNode{
				{
					Detail: "detail message-1",
					Num:    7.7,
				}, {
					Detail: "detail message-2",
					Num:    7.8,
				},
			},
			emptyFault: &[]SimpleNode{},
		},
		{
			name:          "NestedNode",
			wantErrString: "0.00: detail-1\n0.00: nested-2",
			hasData:       true,
			fault: &SimpleNode{
				Detail: "detail-1",
				Num:    .003,
				Nested: &SimpleNode{
					Detail: "nested-2",
					Num:    .004,
				},
			},
			emptyFault: &SimpleNode{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.MarshalIndent(tt.fault, "\t\t\t\t", "\t")
			assert.NoError(t, err)

			var pingRequest = new(Ping)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = xml.NewDecoder(r.Body).Decode(pingRequest)
				rsp := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
										<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
														xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
														xmlns:xsd="http://www.w3.org/2001/XMLSchema">
											<soap:Body>
												<soap:Fault>
													<faultcode>soap:Server</faultcode>
													<faultstring>Custom error message.</faultstring>
													<detail>%v</detail>
												</soap:Fault>
											</soap:Body>
										</soap:Envelope>`, string(data))
				_, _ = w.Write([]byte(rsp))
			}))
			defer ts.Close()

			faultErrString := tt.wantErrString
			client := proxy.NewClient(ts.URL)
			req := &Ping{Request: &PingRequest{Message: "Hi"}}
			fault := Wrapper{Item: tt.emptyFault, hasData: tt.hasData}
			var reply PingResponse
			err = client.CallWithFaultDetail("GetData", req, &reply, &fault)
			assert.EqualError(t, err, faultErrString)
			assert.EqualValues(t, tt.fault, fault.Item)
		})
	}
}

// TestXsdDateTime checks the marshalled xsd datetime
func TestXsdDateTime(t *testing.T) {
	type TestDateTime struct {
		XMLName  xml.Name `xml:"TestDateTime"`
		Datetime *xsd.DateTime
	}
	type TestAttrDateTime struct {
		XMLName  xml.Name      `xml:"TestAttrDateTime"`
		Datetime *xsd.DateTime `xml:"Datetime,attr"`
	}
	{
		// without nanosecond
		testDateTime := TestDateTime{Datetime: xsd.NewDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 0, time.FixedZone("UTC-8", -8*60*60)), true)}
		output, err := xml.MarshalIndent(testDateTime, "", "")
		assert.NoError(t, err)
		expected := "<TestDateTime><Datetime>1951-10-22T01:02:03-08:00</Datetime></TestDateTime>"
		assert.Equal(t, expected, string(output))
	}
	{
		// with nanosecond
		testDateTime := TestDateTime{Datetime: xsd.NewDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.FixedZone("UTC-8", -8*60*60)), true)}
		output, err := xml.MarshalIndent(testDateTime, "", "")
		assert.NoError(t, err)
		expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004-08:00</Datetime></TestDateTime>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling of UTC
	{
		testDateTime := TestDateTime{Datetime: xsd.NewDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.UTC), true)}
		output, err := xml.MarshalIndent(testDateTime, "", "")
		assert.NoError(t, err)
		expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004Z</Datetime></TestDateTime>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling of XsdDateTime without TZ
	{
		testDateTime := TestDateTime{Datetime: xsd.NewDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.UTC), false)}
		output, err := xml.MarshalIndent(testDateTime, "", "")
		assert.NoError(t, err)
		expected := "<TestDateTime><Datetime>1951-10-22T01:02:03.000000004</Datetime></TestDateTime>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling as attribute
	{
		testDateTime := TestAttrDateTime{Datetime: xsd.NewDateTime(time.Date(1951, time.October, 22, 1, 2, 3, 4, time.UTC), true)}
		output, err := xml.MarshalIndent(testDateTime, "", "")
		assert.NoError(t, err)
		expected := "<TestAttrDateTime Datetime=\"1951-10-22T01:02:03.000000004Z\"></TestAttrDateTime>"
		assert.Equal(t, expected, string(output))
	}

	// test unmarshalling
	{
		dateTimes := map[string]time.Time{
			"<TestDateTime><Datetime>1951-10-22T01:02:03.000000004-08:00</Datetime></TestDateTime>": time.Date(1951, time.October, 22, 1, 2, 3, 4, time.FixedZone("-0800", -8*60*60)),
			"<TestDateTime><Datetime>1951-10-22T01:02:03Z</Datetime></TestDateTime>":                time.Date(1951, time.October, 22, 1, 2, 3, 0, time.UTC),
			"<TestDateTime><Datetime>1951-10-22T01:02:03</Datetime></TestDateTime>":                 time.Date(1951, time.October, 22, 1, 2, 3, 0, time.Local),
		}
		for dateTimeStr, dateTimeObj := range dateTimes {
			parsedDt := TestDateTime{}
			err := xml.Unmarshal([]byte(dateTimeStr), &parsedDt)
			assert.NoError(t, err)
			assert.True(t, dateTimeObj.Equal(parsedDt.Datetime.Time()))
		}
	}

	// test unmarshalling as attribute
	{
		dateTimes := map[string]time.Time{
			"<TestAttrDateTime Datetime=\"1951-10-22T01:02:03.000000004-08:00\"></TestAttrDateTime>": time.Date(1951, time.October, 22, 1, 2, 3, 4, time.FixedZone("-0800", -8*60*60)),
			"<TestAttrDateTime Datetime=\"1951-10-22T01:02:03Z\"></TestAttrDateTime>":                time.Date(1951, time.October, 22, 1, 2, 3, 0, time.UTC),
			"<TestAttrDateTime Datetime=\"1951-10-22T01:02:03\"></TestAttrDateTime>":                 time.Date(1951, time.October, 22, 1, 2, 3, 0, time.Local),
		}
		for dateTimeStr, dateTimeObj := range dateTimes {
			parsedDt := TestAttrDateTime{}
			err := xml.Unmarshal([]byte(dateTimeStr), &parsedDt)
			assert.NoError(t, err)
			assert.True(t, dateTimeObj.Equal(parsedDt.Datetime.Time()))
		}
	}
}

// TestXsdDateTime checks the marshalled xsd datetime
func TestXsdDate(t *testing.T) {
	type TestDate struct {
		XMLName xml.Name `xml:"TestDate"`
		Date    *xsd.Date
	}

	type TestAttrDate struct {
		XMLName xml.Name  `xml:"TestAttrDate"`
		Date    *xsd.Date `xml:"Date,attr"`
	}

	// test marshalling
	{
		testDate := TestDate{Date: xsd.NewDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)), false)}
		output, err := xml.MarshalIndent(testDate, "", "")
		assert.NoError(t, err)
		expected := "<TestDate><Date>1951-10-22</Date></TestDate>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling
	{
		testDate := TestDate{Date: xsd.NewDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)), true)}
		output, err := xml.MarshalIndent(testDate, "", "")
		assert.NoError(t, err)
		expected := "<TestDate><Date>1951-10-22-08:00</Date></TestDate>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling of UTC
	{
		testDate := TestDate{
			Date: xsd.NewDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC), true),
		}
		output, err := xml.MarshalIndent(testDate, "", "")
		assert.NoError(t, err)
		expected := "<TestDate><Date>1951-10-22Z</Date></TestDate>"
		assert.Equal(t, expected, string(output))
	}

	// test marshalling as attribute
	{
		testDate := TestAttrDate{Date: xsd.NewDate(time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC), true)}
		output, err := xml.MarshalIndent(testDate, "", "")
		assert.NoError(t, err)
		expected := "<TestAttrDate Date=\"1951-10-22Z\"></TestAttrDate>"
		assert.Equal(t, expected, string(output))
	}

	// test unmarshalling
	{
		dates := map[string]time.Time{
			"<TestDate><Date>1951-10-22</Date></TestDate>":       time.Date(1951, time.October, 22, 0, 0, 0, 0, time.Local),
			"<TestDate><Date>1951-10-22Z</Date></TestDate>":      time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC),
			"<TestDate><Date>1951-10-22-08:00</Date></TestDate>": time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
		}
		for dateStr, dateObj := range dates {
			parsedDt := TestDate{}
			err := xml.Unmarshal([]byte(dateStr), &parsedDt)
			assert.NoError(t, err)
			assert.True(t, dateObj.Equal(parsedDt.Date.Time()))
		}
	}

	// test unmarshalling as attribute
	{
		dates := map[string]time.Time{
			"<TestAttrDate Date=\"1951-10-22\"></TestAttrDate>":       time.Date(1951, time.October, 22, 0, 0, 0, 0, time.Local),
			"<TestAttrDate Date=\"1951-10-22Z\"></TestAttrDate>":      time.Date(1951, time.October, 22, 0, 0, 0, 0, time.UTC),
			"<TestAttrDate Date=\"1951-10-22-08:00\"></TestAttrDate>": time.Date(1951, time.October, 22, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
		}
		for dateStr, dateObj := range dates {
			parsedDate := TestAttrDate{}
			err := xml.Unmarshal([]byte(dateStr), &parsedDate)
			assert.NoError(t, err)
			assert.True(t, dateObj.Equal(parsedDate.Date.Time()))
		}
	}
}

// TestXsdTime checks the marshalled xsd datetime
func TestXsdTime(t *testing.T) {
	type TestTime struct {
		XMLName xml.Name `xml:"TestTime"`
		Time    *xsd.Time
	}

	type TestAttrTime struct {
		XMLName xml.Name  `xml:"TestAttrTime"`
		Time    *xsd.Time `xml:"Time,attr"`
	}

	// test marshalling
	{
		testTime := TestTime{Time: xsd.NewTime(12, 13, 14, 4, time.FixedZone("Test", -19800))}
		output, err := xml.MarshalIndent(testTime, "", "")
		assert.NoError(t, err)
		expected := "<TestTime><Time>12:13:14.000000004-05:30</Time></TestTime>"
		assert.Equal(t, expected, string(output))
	}
	{
		testTime := TestTime{Time: xsd.NewTime(12, 13, 14, 0, time.FixedZone("UTC-8", -8*60*60))}
		output, err := xml.MarshalIndent(testTime, "", "")
		assert.NoError(t, err)
		expected := "<TestTime><Time>12:13:14-08:00</Time></TestTime>"
		assert.Equal(t, expected, string(output))
	}
	{
		testTime := TestTime{Time: xsd.NewTime(12, 13, 14, 0, nil)}
		output, err := xml.MarshalIndent(testTime, "", "")
		assert.NoError(t, err)
		expected := "<TestTime><Time>12:13:14</Time></TestTime>"
		assert.Equal(t, expected, string(output))
	}
	// test marshalling as attribute
	{
		testTime := TestAttrTime{Time: xsd.NewTime(12, 13, 14, 4, time.FixedZone("Test", -19800))}
		output, err := xml.MarshalIndent(testTime, "", "")
		assert.NoError(t, err)
		expected := "<TestAttrTime Time=\"12:13:14.000000004-05:30\"></TestAttrTime>"
		assert.Equal(t, expected, string(output))
	}

	// test unmarshalling without TZ
	{
		timeStr := "<TestTime><Time>12:13:14.000000004</Time></TestTime>"
		parsedTime := TestTime{}
		err := xml.Unmarshal([]byte(timeStr), &parsedTime)
		assert.NoError(t, err)
		assert.Equal(t, 12, parsedTime.Time.Hour())
		assert.Equal(t, 13, parsedTime.Time.Minute())
		assert.Equal(t, 14, parsedTime.Time.Second())
		assert.Equal(t, 4, parsedTime.Time.Nanosecond())
		assert.Nil(t, parsedTime.Time.Location())
	}
	// test unmarshalling with UTC
	{
		timeStr := "<TestTime><Time>12:13:14Z</Time></TestTime>"
		parsedTime := TestTime{}
		err := xml.Unmarshal([]byte(timeStr), &parsedTime)
		assert.NoError(t, err)
		assert.Equal(t, 12, parsedTime.Time.Hour())
		assert.Equal(t, 13, parsedTime.Time.Minute())
		assert.Equal(t, 14, parsedTime.Time.Second())
		assert.Equal(t, 0, parsedTime.Time.Nanosecond())
		assert.Equal(t, "UTC", parsedTime.Time.Location().String())
	}
	// test unmarshalling with non-UTC Tz
	{
		timeStr := "<TestTime><Time>12:13:14-08:00</Time></TestTime>"
		parsedTime := TestTime{}

		err := xml.Unmarshal([]byte(timeStr), &parsedTime)
		assert.NoError(t, err)
		assert.Equal(t, 12, parsedTime.Time.Hour())
		assert.Equal(t, 13, parsedTime.Time.Minute())
		assert.Equal(t, 14, parsedTime.Time.Second())
		assert.Equal(t, 0, parsedTime.Time.Nanosecond())
		_, tzOffset := parsedTime.Time.InnerTime.Zone()
		assert.Equal(t, -8*3600, tzOffset)

	}
	// test unmarshalling as attribute
	{
		timeStr := "<TestAttrTime Time=\"12:13:14Z\"></TestAttrTime>"
		parsedTime := TestAttrTime{}
		err := xml.Unmarshal([]byte(timeStr), &parsedTime)
		assert.NoError(t, err)
		assert.Equal(t, 12, parsedTime.Time.Hour())
		assert.Equal(t, 13, parsedTime.Time.Minute())
		assert.Equal(t, 14, parsedTime.Time.Second())
		assert.Equal(t, 0, parsedTime.Time.Nanosecond())
		assert.Equal(t, "UTC", parsedTime.Time.Location().String())
	}
}

func TestHTTPError(t *testing.T) {
	type httpErrorTest struct {
		name         string
		responseCode int
		responseBody string
		wantErr      bool
		wantErrMsg   string
	}

	tests := []httpErrorTest{
		{
			name:         "should error if server returns 500",
			responseCode: http.StatusInternalServerError,
			responseBody: "internal server error",
			wantErr:      true,
			wantErrMsg:   "HTTP Status 500: internal server error",
		},
		{
			name:         "should error if server returns 403",
			responseCode: http.StatusForbidden,
			responseBody: "forbidden",
			wantErr:      true,
			wantErrMsg:   "HTTP Status 403: forbidden",
		},
		{
			name:         "should not error if server returns 200",
			responseCode: http.StatusOK,
			responseBody: `<?xml version="1.0" encoding="utf-8"?>
							<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
								<soap:Body>
									<PingResponse xmlns="http://example.com/service.xsd">
										<PingResult>
											<Message>Pong hi</Message>
										</PingResult>
									</PingResponse>
								</soap:Body>
							</soap:Envelope>`,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.responseCode)
				_, _ = w.Write([]byte(test.responseBody))
			}))
			defer ts.Close()
			client := proxy.NewClient(ts.URL)
			gotErr := client.Call("GetData", &Ping{}, &PingResponse{})
			if test.wantErr {
				assert.NotNil(t, gotErr)
				requestError, ok := gotErr.(*soap.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, test.responseCode, requestError.StatusCode)
				assert.Equal(t, test.responseBody, string(requestError.ResponseBody))
				assert.Equal(t, test.wantErrMsg, requestError.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
