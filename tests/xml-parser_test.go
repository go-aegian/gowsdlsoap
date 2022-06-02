package tests

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/go-aegian/gowsdlsoap/builder/soap"
	"github.com/go-aegian/gowsdlsoap/proxy"
	"github.com/go-aegian/gowsdlsoap/tests/wsdl-samples/ews/ewsApi"
	"github.com/go-aegian/gowsdlsoap/tests/wsdl-samples/navblue/raidoApi"
	"github.com/stretchr/testify/assert"
)

func TestMarshallXMLNS(t *testing.T) {
	r := soap.NewEnvelope()
	r.AddXmlns("xmlns:tns", "http://test")
	r.Body.Content = raidoApi.Ping{}
	bytes, err := xml.Marshal(&r)
	assert.NoError(t, err)
	fmt.Printf("%s", string(bytes))
}

func TestParseEwsCreateItemResponse(t *testing.T) {
	responseXml := `<?xml version="1.0" encoding="utf-8"?>
				<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
				    <s:Header>
				        <h:ServerVersionInfo MajorVersion="15" MinorVersion="1" MajorBuildNumber="2507" MinorBuildNumber="6" Version="V2017_07_11" xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" xmlns="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
				    </s:Header>
				    <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				        <m:CreateItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
				            <m:ResponseMessages>
				                <m:CreateItemResponseMessage ResponseClass="Success">
				                    <m:ResponseCode>NoError</m:ResponseCode>
				                    <m:Items>
				                        <t:CalendarItem>
				                            <t:ItemId Id="AAMkADU4NWEzY2ExLTI2NGQtNGM1Mi05ZWM1LTllMjhmMjY4ZGMxMABGAAAAAADpOQ9VTAkwSKjfsD+XlxkyBwCQn3OegRRmRp7YcOEXm4lWAAAAAAENAACQn3OegRRmRp7YcOEXm4lWAAAFA+6NAAA=" ChangeKey="DwAAABYAAACQn3OegRRmRp7YcOEXm4lWAAAFA/E2"/>
				                        </t:CalendarItem>
				                    </m:Items>
				                </m:CreateItemResponseMessage>
				            </m:ResponseMessages>
				        </m:CreateItemResponse>
				    </s:Body>
				</s:Envelope>`

	responseObject := soap.EnvelopeResponse{Body: soap.BodyResponse{Content: &ewsApi.CreateItemResponse{}}}

	buffer := strings.NewReader(responseXml)
	dec := xml.NewDecoder(buffer)

	err := dec.Decode(&responseObject)

	assert.NoError(t, err)

	proxy.LogXml("response", responseObject)
}

func TestParseEwsFaultUpdateItemResponse(t *testing.T) {
	responseXml := `<?xml version="1.0" encoding="utf-8"?>
					<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
					    <s:Body>
					        <s:Fault>
					            <faultcode xmlns:a="http://schemas.microsoft.com/exchange/services/2006/types">a:ErrorInvalidArgument</faultcode>
					            <faultstring xml:lang="en-US">The request is invalid.</faultstring>
					            <detail>
					                <e:ResponseCode xmlns:e="http://schemas.microsoft.com/exchange/services/2006/errors">ErrorInvalidArgument</e:ResponseCode>
					                <e:Message xmlns:e="http://schemas.microsoft.com/exchange/services/2006/errors">The request is invalid.</e:Message>
					            </detail>
					        </s:Fault>
					    </s:Body>
					</s:Envelope>`

	responseObject := soap.NewEnvelopeResponse()
	responseObject.Body = soap.BodyResponse{Content: &ewsApi.UpdateItemResponseType{}}

	buffer := strings.NewReader(responseXml)
	dec := xml.NewDecoder(buffer)

	err := dec.Decode(&responseObject)

	assert.NoError(t, err)

	proxy.LogXml("response", responseObject)
}

func TestParseEwsDeleteItemResponse(t *testing.T) {
	responseXml := `<?xml version="1.0" encoding="utf-8"?>
					<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
					    <s:Header>
					        <h:ServerVersionInfo MajorVersion="15" MinorVersion="1" MajorBuildNumber="2507" MinorBuildNumber="6" Version="V2017_07_11" xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" xmlns="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
					    </s:Header>
					    <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
					        <m:DeleteItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
					            <m:ResponseMessages>
					                <m:DeleteItemResponseMessage ResponseClass="Error">
					                    <m:MessageText>The specified object was not found in the store., The process failed to get the correct properties.</m:MessageText>
					                    <m:ResponseCode>ErrorItemNotFound</m:ResponseCode>
					                    <m:DescriptiveLinkKey>0</m:DescriptiveLinkKey>
					                </m:DeleteItemResponseMessage>
					            </m:ResponseMessages>
					        </m:DeleteItemResponse>
					    </s:Body>
					</s:Envelope>`

	responseObject := soap.EnvelopeResponse{Body: soap.BodyResponse{Content: &ewsApi.DeleteItemResponseType{}}}

	buffer := strings.NewReader(responseXml)
	dec := xml.NewDecoder(buffer)

	err := dec.Decode(&responseObject)

	assert.NoError(t, err)

	proxy.LogXml("response", responseObject)
}
