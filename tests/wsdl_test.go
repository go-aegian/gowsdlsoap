package tests

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/go-aegian/gosoap/builder/wsdl"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile(`wsdl-samples\ews\services.wsdl`)
	assert.NoError(t, err)

	v := wsdl.WSDL{}
	err = xml.Unmarshal(data, &v)
	assert.NoError(t, err)
}
