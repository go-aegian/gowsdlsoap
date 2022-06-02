package tests

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/go-aegian/gowsdlsoap/builder/soap"
	"github.com/go-aegian/gowsdlsoap/proxy"
	"github.com/go-aegian/gowsdlsoap/tests/wsdl-samples/ews/ewsApi"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDeleteItemRequest(t *testing.T) {

	eventId := "AAMkADU4NWEzY2ExLTI2NGQtNGM1Mi05ZWM1LTllMjhmMjY4ZGMxMABGAAAAAADpOQ9VTAkwSKjfsD+XlxkyBwCQn3OegRRmRp7YcOEXm4lWAAAAAAENAACQn3OegRRmRp7YcOEXm4lWAAAHwaVOAAA="

	env := soap.NewEnvelope()
	env.Header = &soap.Header{Headers: []interface{}{
		&ewsApi.RequestServerVersion{
			Version: (*ewsApi.ExchangeVersionType)(proxy.String(string(ewsApi.ExchangeVersionTypeExchange2016))),
		}},
	}
	env.Body.Content = &ewsApi.DeleteItemType{
		ItemIds:                  &ewsApi.NonEmptyArrayOfBaseItemIdsType{ItemId: &ewsApi.ItemIdType{Id: eventId}},
		DeleteType:               (*ewsApi.DisposalType)(proxy.String(string(ewsApi.DisposalTypeMoveToDeletedItems))),
		SendMeetingCancellations: (*ewsApi.CalendarItemCreateOrDeleteOperationType)(proxy.String(string(ewsApi.CalendarItemCreateOrDeleteOperationTypeSendToNone))),
	}

	request, err := xml.Marshal(env)
	assert.NoError(t, err)

	fmt.Printf("\nRequest:\n%s\n\n", request)
}

func TestUnmarshallResponse(t *testing.T) {
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
					                            <t:ItemId Id="AAMkADU4NWEzY2ExLTI2NGQtNGM1Mi05ZWM1LTllMjhmMjY4ZGMxMABGAAAAAADpOQ9VTAkwSKjfsD+XlxkyBwCQn3OegRRmRp7YcOEXm4lWAAAAAAENAACQn3OegRRmRp7YcOEXm4lWAAAHwaVOAAA=" ChangeKey="DwAAABYAAACQn3OegRRmRp7YcOEXm4lWAAAHwcYu"/>
					                        </t:CalendarItem>
					                    </m:Items>
					                </m:CreateItemResponseMessage>
					            </m:ResponseMessages>
					        </m:CreateItemResponse>
					    </s:Body>
					</s:Envelope>`

	env := soap.NewEnvelopeResponse()

	env.Body.Content = &ewsApi.CreateItemResponse{}

	err := xml.Unmarshal([]byte(responseXml), env)

	assert.NoError(t, err)

	proxy.LogXml("response", env)
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

	env := soap.NewEnvelopeResponse()
	env.Body.Content = &ewsApi.CreateItemResponse{}

	err := xml.Unmarshal([]byte(responseXml), env)

	assert.NoError(t, err)

	proxy.LogXml("response", env)
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

	env := soap.NewEnvelopeResponse()
	env.Body.Content = &ewsApi.UpdateItemResponseType{}

	err := xml.Unmarshal([]byte(responseXml), env)

	assert.NoError(t, err)

	proxy.LogXml("response", env)
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

	env := soap.NewEnvelopeResponse()
	env.Body.Content = &ewsApi.DeleteItemResponse{}

	err := xml.Unmarshal([]byte(responseXml), env)

	assert.NoError(t, err)

	proxy.LogXml("response", env)
}
