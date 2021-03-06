package tests

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/go-aegian/gowsdlsoap"
	"github.com/stretchr/testify/assert"
)

func TestElementGenerationDoesntCommentOutStructProperty(t *testing.T) {
	g, err := gowsdlsoap.New(`wsdl-samples\test.wsdl`, "soapApi", false, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)

	if strings.Contains(string(resp["types"]), "// this is a comment  GetInfoResult string `xml:\"GetInfoResult,omitempty\"`") {
		t.Error("Type comment should not comment out struct type property")
		t.Error(string(resp["types"]))
	}
}

func TestComplexTypeWithInlineSimpleType(t *testing.T) {
	g, err := gowsdlsoap.New(`wsdl-samples\test.wsdl`, "soapApi", false, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)
	actual, err := getTypeDeclaration(resp, "GetInfo")
	assert.NoError(t, err)

	expected := `type GetInfo struct {
	XMLName	xml.Name	` + "`" + `xml:"http://www.mnb.hu/webservices/ GetInfo"` + "`" + `

	Id	string	` + "`" + `xml:"Id,omitempty" json:"Id,omitempty"` + "`" + `
}`
	if actual != expected {
		t.Error("got " + actual + " want " + expected)
	}
}

//
// func TestAttributeRef(t *testing.T) {
// 	g, err := gowsdlsoap.New(`wsdl-samples\ews\services.wsdl`, "ewsApi", false, true)
// 	assert.NoError(t, err)
//
// 	resp, err := g.Build()
// 	assert.NoError(t, err)
//
// 	actual, err := getTypeDeclaration(resp, "RequestAttachmentIdType")
// 	assert.NoError(t, err)
//
// 	expected := `type RequestAttachmentIdType struct {` + "`" +
// 		`XMLName xml.Name ` + "`" +
// 		`xml:"http://schemas.microsoft.com/exchange/services/2006/types AttachmentId"` + "`" +
// 		` Id string ` + "`" + `xml:"Id,attr,omitempty" json:"Id,omitempty"` + "`" +
// 		`}`
//
// 	actual = strings.TrimSpace(string(bytes.ReplaceAll([]byte(actual), []byte("\n\t"), []byte(" "))))
// 	expected = strings.TrimSpace(string(bytes.ReplaceAll([]byte(expected), []byte("\n\t"), []byte(" "))))
// 	assert.Equal(t, expected, actual)
// }

func TestElementWithLocalSimpleType(t *testing.T) {
	g, err := gowsdlsoap.New(`wsdl-samples\test.wsdl`, "soapApi", false, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)

	// Type declaration
	actual, err := getTypeDeclaration(resp, "ElementWithLocalSimpleType")
	assert.NoError(t, err)

	expected := `type ElementWithLocalSimpleType string`

	assert.Equal(t, expected, actual)

	// Const declaration of first enum value
	actual, err = getTypeDeclaration(resp, "ElementWithLocalSimpleTypeEnum1")
	assert.NoError(t, err)

	expected = `const ElementWithLocalSimpleTypeEnum1 ElementWithLocalSimpleType = "enum1"`

	assert.Equal(t, expected, actual)

	actual, err = getTypeDeclaration(resp, "ElementWithLocalSimpleTypeEnum2")
	assert.NoError(t, err)

	expected = `const ElementWithLocalSimpleTypeEnum2 ElementWithLocalSimpleType = "enum2"`
	assert.Equal(t, expected, actual)
}

func TestDateTimeType(t *testing.T) {
	g, err := gowsdlsoap.New(`wsdl-samples\test.wsdl`, "soapApi", false, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)

	// Type declaration
	actual, err := getTypeDeclaration(resp, "StartDate")
	assert.NoError(t, err)

	expected := `type StartDate xsd.DateTime`

	assert.Equal(t, expected, actual)

	// Method declaration MarshalXML
	actual, err = getFuncDeclaration(resp, "MarshalXML", "StartDate")
	assert.NoError(t, err)

	expected = `func (xdt StartDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return xsd.DateTime(xdt).MarshalXML(e, start)
}`

	assert.Equal(t, expected, actual)

	// Method declaration UnmarshalXML
	actual, err = getFuncDeclaration(resp, "UnmarshalXML", "StartDate")
	assert.NoError(t, err)

	expected = `func (xdt *StartDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return (*xsd.DateTime)(xdt).UnmarshalXML(d, start)
}`

	assert.Equal(t, expected, actual)
}

func TestVboxGeneratesWithoutSyntaxErrors(t *testing.T) {
	files, err := filepath.Glob(`wsdl-samples\*.wsdl`)
	assert.NoError(t, err)

	for _, file := range files {
		g, err := gowsdlsoap.New(file, "soapApi", false, true)
		assert.NoError(t, err)

		resp, err := g.Build()
		if err != nil {
			continue
		}

		data := new(bytes.Buffer)
		data.Write(resp["header"])
		data.Write(resp["types"])
		data.Write(resp["operations"])
		data.Write(resp["soap"])

		_, err = format.Source(data.Bytes())
		assert.NoError(t, err)
	}
}

func TestEnumerationsGeneratedCorrectly(t *testing.T) {
	enumStringTest := func(t *testing.T, fixtureWsdl string, varName string, typeName string, enumString string) {
		g, err := gowsdlsoap.New(`wsdl-samples\`+fixtureWsdl, "soapApi", false, true)
		assert.NoError(t, err)

		resp, err := g.Build()
		assert.NoError(t, err)
		re := regexp.MustCompile(varName + " " + typeName + " = \"([^\"]*)\"")
		matches := re.FindStringSubmatch(string(resp["types"]))

		if len(matches) != 2 {
			t.Errorf("No match or too many matches found for %s", varName)
		} else if matches[1] != enumString {
			t.Errorf("%s got '%s' but expected '%s'", varName, matches[1], enumString)
		}
	}
	enumStringTest(t, "vboxweb.wsdl", "SettingsVersionV1_14", "SettingsVersion", "v1_14")

}

func TestComplexTypeGeneratedCorrectly(t *testing.T) {
	g, err := gowsdlsoap.New(`wsdl-samples\ews\services.wsdl`, "ewsApi", true, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)

	decl, err := getTypeDeclaration(resp, "ItemIdType")

	expected := "type ItemIdType struct"
	re := regexp.MustCompile(expected)
	matches := re.FindStringSubmatch(decl)

	if len(matches) != 1 {
		t.Errorf("No match or too many matches found for ItemIdType")
	} else if matches[0] != expected {
		t.Errorf("ItemIdType got '%s' but expected '%s'", matches[1], expected)
	}
}

func TestEWSWSDL(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	g, err := gowsdlsoap.New(`.\wsdl-samples\ews\services.wsdl`, "ewsApi", true, true)
	assert.NoError(t, err)

	resp, err := g.Build()
	assert.NoError(t, err)
	data := new(bytes.Buffer)
	data.Write(resp["header"])
	data.Write(resp["types"])
	data.Write(resp["operations"])
	data.Write(resp["soap"])

	source, err := format.Source(data.Bytes())
	assert.NoError(t, err)

	if _, err := os.Stat(`.\wsdl-samples\ews\ewsApi\proxy.go`); err != nil {
		_ = ioutil.WriteFile(`.\wsdl-samples\ews\ewsApi\proxy.go`, source, 0664)
	}

	expectedBytes, err := ioutil.ReadFile(`.\wsdl-samples\ews\ewsApi\proxy.go`)
	assert.NoError(t, err)

	actual := string(source)
	expected := string(expectedBytes)
	if actual != expected {
		_ = ioutil.WriteFile(`.\wsdl-samples\ews\ewsApi\proxy_test_gen.go`, source, 0664)
		t.Error(`got source .\wsdl-samples\ews\ewsApi\proxy_test_gen.go but expected .\wsdl-samples\ews\ewsApi\proxy.go`)
	}
}

func getTypeDeclaration(resp map[string][]byte, name string) (string, error) {
	source, err := format.Source([]byte(string(resp["header"]) + string(resp["types"])))
	if err != nil {
		return "", err
	}
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "soapApi.go", string(source), parser.DeclarationErrors)
	if err != nil {
		return "", err
	}
	o := f.Scope.Lookup(name)
	if o == nil {
		return "", errors.New("type " + name + " is missing")
	}
	var buf bytes.Buffer
	buf.WriteString(o.Kind.String())
	buf.WriteString(" ")
	err = printer.Fprint(&buf, fileSet, o.Decl)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func findFuncDecl(f *ast.File, name string, recv string) *ast.Decl {
	// Loop over all declarations
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Found FuncDecl declaration type
			if funcDecl.Name.Name == name {
				// Found match with function name
				if ident, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
					// Value receiver type
					if ident.Name == recv {
						// Found receiver type match
						return &decl
					}
				} else if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					// Pointer receiver type
					if t, ok := starExpr.X.(*ast.Ident); ok {
						if t.Name == recv {
							// Found receiver type match
							return &decl
						}
					}
				}
			}
		}
	}

	return nil
}

func getFuncDeclaration(resp map[string][]byte, name string, recv string) (string, error) {
	source, err := format.Source([]byte(string(resp["header"]) + string(resp["types"])))
	if err != nil {
		return "", err
	}
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "soapApi.go", string(source), parser.DeclarationErrors)
	if err != nil {
		return "", err
	}

	decl := findFuncDecl(f, name, recv)
	if decl == nil {
		return "", errors.New("Function declaration " + name + " not found")
	}
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fileSet, *decl)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
