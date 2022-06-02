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
	}, nil
}

// Build initiates the code generation process by starting two goroutines:
//   generate types
//   generate operations
func (b *Builder) Build() (map[string][]byte, error) {
	code := make(map[string][]byte)

	err := b.unmarshal()
	if err != nil {
		return nil, err
	}

	// Process WSDL nodes
	for _, schema := range b.wsdl.Types.Schemas {
		NewXsdParser(schema, b.wsdl.Types.Schemas).parse()
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		code["types"], err = b.parseTypes()
		if err != nil {
			log.Println("parseTypes", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		code["operations"], err = b.parseOperations()
		if err != nil {
			log.Println("parseOperations", "error", err)
		}
	}()

	wg.Wait()

	code["header"], err = b.parseHeader()
	if err != nil {
		log.Println(err)
	}

	return code, nil
}

// Method setNamespace sets (and returns) the currently active XML namespace.
func (b *Builder) setNamespace(ns string) string {
	b.currentNamespace = ns
	return b.currentNamespace
}

// Method setNamespace returns the currently active XML namespace.
func (b *Builder) getNamespace() string {
	return b.currentNamespace
}

func (b *Builder) readFile(loc *location) (data []byte, err error) {
	if loc.file != "" {
		log.Println("Reading", "file", loc.file)
		data, err = ioutil.ReadFile(loc.file)
		return
	}

	log.Println("Downloading", "file", loc.url.String())
	data, err = downloadFile(loc.url.String(), b.skipTls)
	return
}

func (b *Builder) unmarshal() error {
	data, err := b.readFile(b.location)
	if err != nil {
		return err
	}

	b.wsdl = new(wsdl.WSDL)
	err = xml.Unmarshal(data, b.wsdl)
	if err != nil {
		return err
	}

	for _, schema := range b.wsdl.Types.Schemas {
		err = b.resolveExternal(schema, b.location)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) resolveExternal(schema *xsd.Schema, loc *location) error {
	download := func(base *location, ref string) error {
		location, err := base.Parse(ref)
		if err != nil {
			return err
		}

		schemaKey := location.String()
		if b.xsdExternals[location.String()] {
			return nil
		}

		if b.xsdExternals == nil {
			b.xsdExternals = make(map[string]bool, maxRecursion)
		}

		b.xsdExternals[schemaKey] = true

		var data []byte
		if data, err = b.readFile(location); err != nil {
			return err
		}

		newSchema := new(xsd.Schema)

		err = xml.Unmarshal(data, newSchema)
		if err != nil {
			return err
		}

		if (len(newSchema.Includes) > 0 || len(newSchema.Imports) > 0) && maxRecursion > b.currentRecursionLevel {
			b.currentRecursionLevel++

			err = b.resolveExternal(newSchema, location)
			if err != nil {
				return err
			}
		}

		b.wsdl.Types.Schemas = append(b.wsdl.Types.Schemas, newSchema)

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

func (b *Builder) parseTypes() ([]byte, error) {
	funcMap := template.FuncMap{
		"isBasicType":              isBasicType,
		"toGoType":                 toGoType,
		"stripAliasNSFromType":     stripAliasNSFromType,
		"replaceReservedWords":     replaceReservedWords,
		"replaceAttrReservedWords": replaceAttrReservedWords,
		"normalize":                normalize,
		"makeFieldPublic":          makePublic,
		"comment":                  comment,
		"goString":                 goString,
		"isInnerBasicType":         b.isInnerBasicType,
		"isAbstract":               b.isAbstract,
		"makePublic":               b.makePublicFn,
		"findMessageType":          b.findMessageType,
		"findNameByType":           b.findNameByType,
		"stripPointerFromType":     stripPointerFromType,
		"setNamespace":             b.setNamespace,
		"getNamespace":             b.getNamespace,
		"packageName":              b.packageName,
		"getAliasNS":               getAliasNS,
		"getNSFromType":            b.getNSFromType,
		"getNSAlias":               b.getNSAlias,
		"getNS":                    getNS,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("types").Funcs(funcMap).Parse(templates.Types))

	err := tmpl.Execute(data, b.wsdl.Types)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Builder) parseOperations() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":             toGoType,
		"stripAliasNSFromType": stripAliasNSFromType,
		"replaceReservedWords": replaceReservedWords,
		"normalize":            normalize,
		"makePrivate":          makePrivate,
		"packageName":          b.packageName,
		"makePublic":           b.makePublicFn,
		"findMessageType":      b.findMessageType,
		"findSOAPAction":       b.findSOAPAction,
		"findServiceAddress":   b.findServiceAddress,
		"getXmlns":             b.getXmlns,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("operations").Funcs(funcMap).Parse(templates.Operations))

	err := tmpl.Execute(data, b.wsdl.PortTypes)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Builder) parseHeader() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":             toGoType,
		"stripAliasNSFromType": stripAliasNSFromType,
		"replaceReservedWords": replaceReservedWords,
		"normalize":            normalize,
		"comment":              comment,
		"makePublic":           b.makePublicFn,
		"findMessageType":      b.findMessageType,
	}

	data := new(bytes.Buffer)

	tmpl := template.Must(template.New("header").Funcs(funcMap).Parse(templates.Header))

	err := tmpl.Execute(data, b.pkg)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (b *Builder) packageName() string {
	return b.pkg
}

func (b *Builder) isAbstract(t string, checkParent bool) bool {
	t = stripAliasNSFromType(t)
	if isBasicType(t) {
		return true
	}

	for _, schema := range b.wsdl.Types.Schemas {
		for _, complexType := range schema.ComplexTypes {
			if complexType.Name == t {
				if checkParent {
					if complexType.Abstract {
						return true
					}
					baseType := stripAliasNSFromType(complexType.ComplexContent.Extension.Base)

					if baseType == "" {
						return false
					}

					for _, complexTypeInner := range schema.ComplexTypes {
						if complexTypeInner.Name == baseType && complexTypeInner.Abstract {
							return true
						}
					}
				} else {
					return complexType.Abstract
				}
			}
		}
	}

	return false
}

func (b *Builder) isInnerBasicType(t string) bool {
	t = stripAliasNSFromType(t)
	if isBasicType(t) {
		return true
	}

	for _, schema := range b.wsdl.Types.Schemas {
		for _, simpleType := range schema.SimpleType {
			if simpleType.Name == t {
				return true
			}
		}
	}

	for _, schema := range b.wsdl.Types.Schemas {
		for _, complexType := range schema.ComplexTypes {
			if complexType.Name == t && !complexType.Mixed && (len(complexType.Sequence) > 0 || len(complexType.Choice) > 0 || len(complexType.SequenceChoice) > 0 || complexType.Abstract) {
				return true
			}
		}
	}

	return false
}

func (b *Builder) findMessageType(message string) string {
	message = stripAliasNSFromType(message)

	for _, msg := range b.wsdl.Messages {
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
			return stripAliasNSFromType(part.Type)
		}

		elRef := stripAliasNSFromType(part.Element)

		for _, schema := range b.wsdl.Types.Schemas {
			for _, el := range schema.Elements {
				if strings.EqualFold(elRef, el.Name) {
					if el.Type != "" {
						return stripAliasNSFromType(el.Type)
					}

					return el.Name
				}
			}
		}
	}
	return ""
}

func (b *Builder) getNSAlias(ns string) string {
	for _, schema := range b.wsdl.Types.Schemas {
		if schema.TargetNamespace == ns {
			for alias, url := range schema.Xmlns {
				if url == ns {
					return alias + ":"
				}
			}
		}
	}
	return ""
}

func (b *Builder) getNSFromType(ns string) string {
	aliasNS := getAliasNS(ns)

	for _, schema := range b.wsdl.Types.Schemas {
		for alias, url := range schema.Xmlns {
			if alias == aliasNS {
				return url
			}
		}
	}
	return ""
}

func getAliasNS(typeName string) string {
	r := strings.Split(typeName, ":")
	if len(r) == 2 && r[0] != "xs" {
		return r[0] + ":"
	}
	return ""
}

func getNS(ns string) string {
	r := strings.Split(ns, ":")
	if len(r) == 1 {
		return r[0]
	}
	return r[1]
}

// Given a type, check if there's an Element with that type, and return its name.
func (b *Builder) findNameByType(name string, getNS bool) string {
	return NewXsdParser(nil, b.wsdl.Types.Schemas).findNameByType(name, getNS)
}

func (b *Builder) findSOAPAction(operation, portType string) string {
	for _, binding := range b.wsdl.Binding {
		if strings.ToUpper(stripAliasNSFromType(binding.Type)) != strings.ToUpper(portType) {
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

func (b *Builder) findServiceAddress(name string) string {
	for _, service := range b.wsdl.Service {
		for _, port := range service.Ports {
			if port.Name == name {
				return port.SOAPAddress.Location
			}
		}
	}

	return ""
}

func (b *Builder) getXmlns() map[string]string {
	for alias, url := range b.wsdl.Xmlns {
		if alias == "tns" {
			for _, schema := range b.wsdl.Types.Schemas {
				if schema.TargetNamespace == url {
					return schema.Xmlns
				}
			}
		}
	}
	return map[string]string{}
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

func toGoType(xsdType string, nillable bool) string {
	t := stripAliasNSFromType(xsdType)
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

func stripAliasNSFromType(fullType string) string {
	r := strings.Split(fullType, ":")
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
	_, exists := basicTypes[stripAliasNSFromType(identifier)]
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
