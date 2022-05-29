package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-aegian/gosoap/builder"
)

func TestLocation_ParseLocation_URL(t *testing.T) {
	r, err := builder.NewLocation("http://example.org/my.wsdl")
	if err != nil {
		t.Fatal(err)
	}

	if !r.IsURL() || r.IsFile() {
		t.Error("location should be a URL type")
	}
	if r.String() != "http://example.org/my.wsdl" {
		t.Error("got " + r.String() + " wanted " + "http://example.org/my.wsdl")
	}
}

func TestLocation_Parse_URL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"http://example.org/my.wsdl", "some.xsd", "http://example.org/some.xsd"},
		{"http://example.org/folder/my.wsdl", "some.xsd", "http://example.org/folder/some.xsd"},
		{"http://example.org/folder/my.wsdl", "../some.xsd", "http://example.org/some.xsd"},
	}
	for _, test := range tests {
		r, err := builder.NewLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if !r.IsURL() || r.IsFile() {
			t.Error("location should be a URL type")
		}
		if r.String() != test.expected {
			t.Error("got " + r.String() + " wanted " + test.name)
		}
	}
}

func TestLocation_ParseLocation_File(t *testing.T) {
	tests := []struct {
		name string
	}{
		{`wsdl-samples\test.wsdl`},
		{`wsdl-samples\test.wsdl`},
	}
	for _, test := range tests {
		r, err := builder.NewLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}

		if r.IsURL() || !r.IsFile() {
			t.Error("location should be a FILE type")
			continue
		}
		if !filepath.IsAbs(r.String()) {
			t.Error("Path should be absolute")
		}
		if _, err := os.Stat(r.String()); err != nil {
			t.Errorf("location should point to existing loc: %s", err.Error())
		}
	}
}

func TestLocation_Parse_File(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{`wsdl-samples\test.wsdl`, `some.xsd`, `wsdl-samples\some.xsd`},
		{`wsdl-samples\test.wsdl`, `..\xsd\some.xsd`, `xsd\some.xsd`},
		{`wsdl-samples\test.wsdl`, `xsd\some.xsd`, `wsdl-samples\xsd\some.xsd`},
	}
	for _, test := range tests {
		r, err := builder.NewLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if r.IsURL() || !r.IsFile() {
			t.Error("location should be a File type")
			continue
		}
		x, _ := filepath.Abs("")
		rel, _ := filepath.Rel(x, r.String())
		if rel != test.expected {
			t.Error("got " + rel + " wanted " + test.expected)
		}
	}
}

func TestLocation_Parse_FileToURL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{`wsdl-samples\test.wsdl`, "http://example.org/some.xsd", "http://example.org/some.xsd"},
	}
	for _, test := range tests {
		r, err := builder.NewLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if !r.IsURL() || r.IsFile() {
			t.Error("location should be a URL type")
			continue
		}
		if r.String() != test.expected {
			t.Error("got " + r.String() + " wanted " + test.expected)
		}
	}
}
