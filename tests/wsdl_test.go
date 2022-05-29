package tests

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/go-aegian/gosoap/builder/wsdl"
)

func TestUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile(`wsdl-samples\ews\services.wsdl`)
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}

	v := wsdl.WSDL{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}
}
