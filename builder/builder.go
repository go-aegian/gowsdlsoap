package builder

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"

	"github.com/go-aegian/gowsdlsoap/builder/templates"
	"github.com/go-aegian/gowsdlsoap/builder/wsdl"
	"github.com/go-aegian/gowsdlsoap/builder/xsd"
)

const maxRecursion uint8 = 20

var basicTypes = map[string]string{
	"string":      "string",
	"float32":     "float32",
	"float64":     "float64",
	"int":         "int",
	"int8":        "int8",
	"int16":       "int16",
	"int32":       "int32",
	"int64":       "int64",
	"bool":        "bool",
	"time.Time":   "time.Time",
	"[]byte":      "[]byte",
	"byte":        "byte",
	"uint":        "uint",
	"uint8":       "uint8",
	"uint16":      "uint16",
	"uint32":      "uint32",
	"uint64":      "uint64",
	"interface{}": "interface{}",
}

var xsd2GoTypes = map[string]string{
	"string":        "string",
	"token":         "string",
	"float":         "float32",
	"double":        "float64",
	"decimal":       "float64",
	"integer":       "int32",
	"int":           "int32",
	"short":         "int16",
	"byte":          "int8",
	"long":          "int64",
	"boolean":       "bool",
	"datetime":      "xsd.DateTime",
	"date":          "xsd.Date",
	"time":          "xsd.Time",
	"base64binary":  "[]byte",
	"hexbinary":     "[]byte",
	"unsignedint":   "uint32",
	"unsignedshort": "uint16",
	"unsignedbyte":  "byte",
	"unsignedlong":  "uint64",
	"anytype":       "AnyType",
	"ncname":        "NCName",
	"anyuri":        "AnyURI",
}

var reservedWords = map[string]string{
	"break":       "break_",
	"default":     "default_",
	"func":        "func_",
	"interface":   "interface_",
	"select":      "select_",
	"case":        "case_",
	"defer":       "defer_",
	"go":          "go_",
	"map":         "map_",
	"struct":      "struct_",
	"chan":        "chan_",
	"else":        "else_",
	"goto":        "goto_",
	"package":     "package_",
	"switch":      "switch_",
	"const":       "const_",
	"fallthrough": "fallthrough_",
	"if":          "if_",
	"range":       "range_",
	"type":        "type_",
	"continue":    "continue_",
	"for":         "for_",
	"import":      "import_",
	"return":      "return_",
	"var":         "var_",
}

var reservedWordsInAttr = map[string]string{
	"break":       "break_",
	"default":     "default_",
	"func":        "func_",
	"interface":   "interface_",
	"select":      "select_",
	"case":        "case_",
	"defer":       "defer_",
	"go":          "go_",
	"map":         "map_",
	"struct":      "struct_",
	"chan":        "chan_",
	"else":        "else_",
	"goto":        "goto_",
	"package":     "package_",
	"switch":      "switch_",
	"const":       "const_",
	"fallthrough": "fallthrough_",
	"if":          "if_",
	"range":       "range_",
	"type":        "type_",
	"continue":    "continue_",
	"for":         "for_",
	"import":      "import_",
	"return":      "return_",
	"var":         "var_",
	"string":      "string_",
}

var timeout = 30 * time.Second

var cacheDir = filepath.Join(os.TempDir(), "gowsdlsoap-cache")

func init() {
	err := os.MkdirAll(cacheDir, 0700)
	if err != nil {
		log.Println("create cache directory", "error", err)
		os.Exit(1)
	}
}

// Builder defines the struct for WSDL generator.
type Builder struct {
	location              *location
	pkg                   string
	skipTls               bool
	makePublicFn          func(string) string
	wsdl                  *wsdl.WSDL
	xsdExternals          map[string]bool
	currentRecursionLevel uint8
	currentNamespace      string
	hasXMLName            bool
}

func New(file, pkg string, ignoreTLS bool, exportAllTypes bool) (*Builder, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		return nil, errors.New("WSDL file is required to generate Go proxy")
	}

	pkg = strings.TrimSpace(pkg)
	if pkg == "" {
		pkg = "soapProxy"
	}

	makePublicFn := func(id string) string { return id }
	if exportAllTypes {
		makePublicFn = makePublic
	}

	r, err := NewLocation(file)
	if err != nil {
		return nil, err
	}

	return &Builder{
		location:     r,
		pkg:          pkg,
		skipTls:      ignoreTLS,
		makePublicFn: makePublicFn,
		hasXMLName:   false,
	}, nil
}

// Build initiates the code generation process by starting two goroutines:
//   generate types
//   generate operations
func (g *Builder) Build() (map[string][]byte, error) {
	code := make(map[string][]byte)

	err := g.unmarshal()
	if err != nil {
		return nil, err
	}

	// Process WSDL nodes
	for _, schema := range g.wsdl.Types.Schemas {
		NewXsdParser(schema, g.wsdl.Types.Schemas).parse()
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		code["types"], err = g.parseTypes()
		if err != nil {
			log.Println("parseTypes", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		code["operations"], err = g.parseOperations()
		if err != nil {
			log.Println("parseOperations", "error", err)
		}
	}()

	wg.Wait()

	code["header"], err = g.parseHeader()
	if err != nil {
		log.Println(err)
	}

	return code, nil
}

// Method setNamespace sets (and returns) the currently active XML namespace.
func (g *Builder) setNamespace(ns string) string {
	g.currentNamespace = ns
	return g.currentNamespace
}

// Method setNamespace returns the currently active XML namespace.
func (g *Builder) getNamespace() string {
	return g.currentNamespace
}

func (g *Builder) readFile(loc *location) (data []byte, err error) {
	if loc.file != "" {
		log.Println("Reading", "file", loc.file)
		data, err = ioutil.ReadFile(loc.file)
		return
	}

	log.Println("Downloading", "file", loc.url.String())
	data, err = downloadFile(loc.url.String(), g.skipTls)
	return
}

func (g *Builder) unmarshal() error {
	data, err := g.readFile(g.location)
	if err != nil {
		return err
	}

	g.wsdl = new(wsdl.WSDL)
	err = xml.Unmarshal(data, g.wsdl)
	if err != nil {
		return err
	}

	for _, schema := range g.wsdl.Types.Schemas {
		err = g.resolveExternal(schema, g.location)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Builder) resolveExternal(schema *xsd.Schema, loc *location) error {
	download := func(base *location, ref string) error {
		location, err := base.Parse(ref)
		if err != nil {
			return err
		}

		schemaKey := location.String()
		if g.xsdExternals[location.String()] {
			return nil
		}

		if g.xsdExternals == nil {
			g.xsdExternals = make(map[string]bool, maxRecursion)
		}

		g.xsdExternals[schemaKey] = true

		var data []byte
		if data, err = g.readFile(location); err != nil {
			return err
		}

		newSchema := new(xsd.Schema)

		err = xml.Unmarshal(data, newSchema)
		if err != nil {
			return err
		}

		if (len(newSchema.Includes) > 0 || len(newSchema.Imports) > 0) &&
			maxRecursion > g.currentRecursionLevel {
			g.currentRecursionLevel++

			err = g.resolveExternal(newSchema, location)
			if err != nil {
				return err
			}
		}

		g.wsdl.Types.Schemas = append(g.wsdl.Types.Schemas, newSchema)

		return nil
	}

	for _, xsdImport := range schema.Imports {
		// Download the file only if we have a hint in the form of schemaLocation.
		if xsdImport.SchemaLocation == "" {
			log.Printf("[WARN] Don't know where to find XSD for %s", xsdImport.Namespace)
			continue
		}

		if e := download(loc, xsdImport.SchemaLocation); e != nil {
			return e
		}
	}

	for _, incl := range schema.Includes {
		if e := download(loc, incl.SchemaLocation); e != nil {
			return e
		}
	}

	return nil
}

func (g *Builder) parseTypes() ([]byte, error) {
	funcMap := template.FuncMap{
		"isBasicType":              isBasicType,
		"toGoType":                 toGoType,
		"stripNamespaceFromType":   stripNamespaceFromType,
		"replaceReservedWords":     replaceReservedWords,
		"replaceAttrReservedWords": replaceAttrReservedWords,
		"normalize":                normalize,
		"makeFieldPublic":          makePublic,
		"comment":                  comment,
		"stripNamespace":           stripNamespace,
		"goString":                 goString,
		"setHasXMLName":            g.setHasXMLName,
		"getHasXMLName":            g.getHasXMLName,
		"isInnerBasicType":         g.isInnerBasicType,
		"isAbstract":               g.isAbstract,
		"makePublic":               g.makePublicFn,
		"findMessageType":          g.findMessageType,
		"findNameByType":           g.findNameByType,
		"stripPointerFromType":     stripPointerFromType,
		"setNamespace":             g.setNamespace,
		"getNamespace":             g.getNamespace,
		"packageName":              g.packageName,
		"outputNSInField":          g.outputNSInField,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("types").Funcs(funcMap).Parse(templates.Types))

	err := tmpl.Execute(data, g.wsdl.Types)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (g *Builder) parseOperations() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":               toGoType,
		"stripNamespaceFromType": stripNamespaceFromType,
		"replaceReservedWords":   replaceReservedWords,
		"normalize":              normalize,
		"makePrivate":            makePrivate,
		"packageName":            g.packageName,
		"makePublic":             g.makePublicFn,
		"findMessageType":        g.findMessageType,
		"findSOAPAction":         g.findSOAPAction,
		"findServiceAddress":     g.findServiceAddress,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("operations").Funcs(funcMap).Parse(templates.Operations))

	err := tmpl.Execute(data, g.wsdl.PortTypes)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (g *Builder) parseHeader() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":               toGoType,
		"stripNamespaceFromType": stripNamespaceFromType,
		"replaceReservedWords":   replaceReservedWords,
		"normalize":              normalize,
		"comment":                comment,
		"makePublic":             g.makePublicFn,
		"findMessageType":        g.findMessageType,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("header").Funcs(funcMap).Parse(templates.Header))

	err := tmpl.Execute(data, g.pkg)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (g *Builder) packageName() string {
	return g.pkg
}

func (g *Builder) isAbstract(t string) bool {
	t = stripNamespaceFromType(t)
	if isBasicType(t) {
		return true
	}

	for _, schema := range g.wsdl.Types.Schemas {
		for _, complexType := range schema.ComplexTypes {
			if complexType.Name == t {
				if complexType.Abstract {
					return true
				}

				baseType := stripNamespaceFromType(complexType.ComplexContent.Extension.Base)

				if baseType == "" {
					return false
				}

				for _, complexTypeInner := range schema.ComplexTypes {
					if complexTypeInner.Name == baseType && complexTypeInner.Abstract {
						return true
					}
				}
			}
		}
	}

	return false
}

func (g *Builder) outputNSInField(t string) bool {
	t = stripNamespaceFromType(t)
	if isBasicType(t) {
		return true
	}

	for _, schema := range g.wsdl.Types.Schemas {
		for _, simpleType := range schema.SimpleType {
			if simpleType.Name == t {
				return true
			}
		}

		for _, complexType := range schema.ComplexTypes {
			if complexType.Name == t && (len(complexType.Sequence) > 0 || len(complexType.Choice) > 0 || len(complexType.SequenceChoice) > 0) && !complexType.Abstract {
				return false
			}
		}
	}

	return false
}

func (g *Builder) setHasXMLName(b bool) bool {
	g.hasXMLName = b
	return g.hasXMLName
}
func (g *Builder) getHasXMLName() bool {
	return g.hasXMLName
}

func (g *Builder) isInnerBasicType(t string) bool {
	t = stripNamespaceFromType(t)
	if isBasicType(t) {
		return true
	}

	for _, schema := range g.wsdl.Types.Schemas {
		for _, simpleType := range schema.SimpleType {
			if simpleType.Name == t {
				return true
			}
		}

		for _, complexType := range schema.ComplexTypes {
			if complexType.Name == t && !complexType.Mixed && (len(complexType.Sequence) > 0 || len(complexType.Choice) > 0 || len(complexType.SequenceChoice) > 0 || complexType.Abstract) {
				return true
			}
		}
	}

	return false
}

func (g *Builder) findMessageType(message string) string {
	message = stripNamespaceFromType(message)

	for _, msg := range g.wsdl.Messages {
		if msg.Name != message {
			continue
		}

		// Assumes document/literal wrapped WS-I
		if len(msg.Parts) == 0 {
			// Message does not have parts.
			// This could be a Port with HTTP binding or SOAP 1.2 binding, which are not currently supported.
			log.Printf("[WARN] %s message doesn't have any parts, ignoring message...", msg.Name)
			continue
		}

		part := msg.Parts[0]
		if part.Type != "" {
			return stripNamespaceFromType(part.Type)
		}

		elRef := stripNamespaceFromType(part.Element)

		for _, schema := range g.wsdl.Types.Schemas {
			for _, el := range schema.Elements {
				if strings.EqualFold(elRef, el.Name) {
					if el.Type != "" {
						return stripNamespaceFromType(el.Type)
					}

					return el.Name
				}
			}
		}
	}
	return ""
}

// Given a type, check if there's an Element with that type, and return its name.
func (g *Builder) findNameByType(name string) string {
	return NewXsdParser(nil, g.wsdl.Types.Schemas).findNameByType(name)
}

func (g *Builder) findSOAPAction(operation, portType string) string {
	for _, binding := range g.wsdl.Binding {
		if strings.ToUpper(stripNamespaceFromType(binding.Type)) != strings.ToUpper(portType) {
			continue
		}

		for _, soapOp := range binding.Operations {
			if soapOp.Name == operation {
				return soapOp.SOAPOperation.SOAPAction
			}
		}
	}

	return ""
}

func (g *Builder) findServiceAddress(name string) string {
	for _, service := range g.wsdl.Service {
		for _, port := range service.Ports {
			if port.Name == name {
				return port.SOAPAddress.Location
			}
		}
	}

	return ""
}

// replaceReservedWords Go reserved keywords to avoid compilation issues
func replaceReservedWords(identifier string) string {
	value := reservedWords[identifier]
	if value != "" {
		return value
	}

	return normalize(identifier)
}

// replaceAttrReservedWords Go reserved keywords to avoid compilation issues
func replaceAttrReservedWords(identifier string) string {
	value := reservedWordsInAttr[identifier]
	if value != "" {
		return value
	}

	return normalize(identifier)
}

// Normalizes value to be used as a valid Go identifier, avoiding compilation issues
func normalize(value string) string {
	mapping := func(r rune) rune {
		if r == '.' {
			return '_'
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}

		return -1
	}

	return strings.Map(mapping, value)
}

func goString(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}

func stripNamespace(xsdType string) string {
	// Handles name space, ie. xsd:string, xs:string
	r := strings.Split(xsdType, ":")

	if len(r) == 2 {
		return r[1]
	}

	return r[0]
}

func toGoType(xsdType string, nillable bool) string {
	// Handles name space, ie. xsd:string, xs:string
	r := strings.Split(xsdType, ":")
	t := r[0]

	if len(r) == 2 {
		t = r[1]
	}

	value := xsd2GoTypes[strings.ToLower(t)]

	if value != "" {
		if nillable {
			value = "*" + value
		}
		return value
	}

	return "*" + replaceReservedWords(makePublic(t))
}

func stripPointerFromType(goType string) string {
	return regexp.MustCompile("^\\s*\\*").ReplaceAllLiteralString(goType, "")
}

func stripNamespaceFromType(xsdType string) string {
	r := strings.Split(xsdType, ":")
	t := r[0]

	if len(r) == 2 {
		t = r[1]
	}

	return strings.Trim(t, "*")
}

func makePublic(identifier string) string {
	if isBasicType(identifier) {
		return identifier
	}

	if identifier == "" {
		return "EmptyString"
	}

	field := []rune(identifier)
	if len(field) == 0 {
		return identifier
	}

	field[0] = unicode.ToUpper(field[0])
	return string(field)
}

func makePrivate(identifier string) string {
	field := []rune(identifier)
	if len(field) == 0 {
		return identifier
	}

	field[0] = unicode.ToLower(field[0])
	return string(field)
}

func isBasicType(identifier string) bool {

	_, exists := basicTypes[stripNamespaceFromType(identifier)]
	return exists
}

func comment(text string) string {
	lines := strings.Split(text, "\n")

	var output string
	if len(lines) == 1 && lines[0] == "" {
		return ""
	}

	hasComment := false

	for _, line := range lines {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
		if line == "" {
			continue
		}
		hasComment = true
		output += "\n// " + line
	}

	if hasComment {
		return output
	}

	return ""
}

func downloadFile(url string, ignoreTLS bool) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreTLS},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := net.Dialer{Timeout: timeout}
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received response code %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
