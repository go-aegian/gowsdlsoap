package builder

import (
	"encoding/xml"
	"strings"

	"github.com/go-aegian/gosoap/builder/xsd"
)

type parseMode int32

const (
	refResolution parseMode = iota
	findNameByType
)

type xsdParser struct {
	c                 *xsd.Schema
	all               []*xsd.Schema
	mode              parseMode
	typeName          string
	foundElementName  string
	typeUsageConflict bool
}

func NewXsdParser(c *xsd.Schema, all []*xsd.Schema) *xsdParser {
	return &xsdParser{
		c:    c,
		all:  all,
		mode: refResolution,
	}
}

func (t *xsdParser) parse() {
	t.mode = refResolution

	for _, ct := range t.c.ComplexTypes {
		t.parseComplexType(ct)
	}

	for _, st := range t.c.SimpleType {
		t.parseSimpleType(st)
	}

	for _, elm := range t.c.Elements {
		t.parseElement(elm)
	}
}

func (t *xsdParser) parseElements(elements []*xsd.Element) {
	for _, elm := range elements {
		t.parseElement(elm)
	}
}

func (t *xsdParser) parseElement(element *xsd.Element) {
	t.findElementName(element)

	if element.ComplexType != nil {
		t.parseComplexType(element.ComplexType)
	}

	if element.SimpleType != nil {
		t.parseSimpleType(element.SimpleType)
	}
}

func (t *xsdParser) parseSimpleType(st *xsd.SimpleType) {
}

func (t *xsdParser) parseComplexType(ct *xsd.ComplexType) {
	t.parseElements(ct.Sequence)
	t.parseElements(ct.Choice)
	t.parseElements(ct.SequenceChoice)
	t.parseElements(ct.All)
	t.parseAttributes(ct.Attributes)
	t.parseAttributes(ct.ComplexContent.Extension.Attributes)
	t.parseElements(ct.ComplexContent.Extension.Sequence)
	t.parseElements(ct.ComplexContent.Extension.Choice)
	t.parseElements(ct.ComplexContent.Extension.SequenceChoice)
	t.parseAttributes(ct.SimpleContent.Extension.Attributes)
}

func (t *xsdParser) parseAttributes(attrs []*xsd.Attribute) {
	for _, attr := range attrs {
		t.parseAttribute(attr)
	}
}

func (t *xsdParser) parseAttribute(attr *xsd.Attribute) {
	if t.mode != refResolution {
		return
	}

	if attr.Ref != "" {
		refAttr := t.getGlobalAttribute(attr.Ref)
		if refAttr != nil && refAttr.Ref == "" {
			t.parseAttribute(refAttr)
			attr.Name = refAttr.Name
			attr.Type = refAttr.Type
			if attr.Fixed == "" {
				attr.Fixed = refAttr.Fixed
			}
		}
	} else if attr.Type == "" {
		if attr.SimpleType != nil {
			t.parseSimpleType(attr.SimpleType)
			attr.Type = attr.SimpleType.Restriction.Base
		}
	}
}

func (t *xsdParser) findElementName(element *xsd.Element) {
	if t.mode != findNameByType {
		return
	}

	if t.typeUsageConflict {
		return
	}

	if stripNamespaceFromType(element.Type) == t.typeName {
		if len(t.foundElementName) == 0 {
			t.foundElementName = element.Name
		} else if t.foundElementName != element.Name {
			// Duplicate use of t.typeName with different element names
			t.typeUsageConflict = true
		}
	}
}

func (t *xsdParser) getGlobalAttribute(name string) *xsd.Attribute {
	ref := t.buildQualifiedName(name)

	for _, schema := range t.all {
		if schema.TargetNamespace == ref.Space {
			for _, attr := range schema.Attributes {
				if attr.Name == ref.Local {
					return attr
				}
			}
		}
	}

	return nil
}

func (t *xsdParser) buildQualifiedName(name string) (qualifiedName xml.Name) {
	x := strings.SplitN(name, ":", 2)
	if len(x) == 1 {
		qualifiedName.Local = x[0]
	} else {
		qualifiedName.Local = x[1]
		qualifiedName.Space = x[0]
		if ns, ok := t.c.Xmlns[qualifiedName.Space]; ok {
			qualifiedName.Space = ns
		}
	}

	return qualifiedName
}

func (t *xsdParser) initFindNameByType(name string) {
	t.mode = findNameByType
	t.typeName = stripNamespaceFromType(name)
	t.foundElementName = ""
	t.typeUsageConflict = false
}

func (t *xsdParser) findNameByType(name string) string {
	t.initFindNameByType(name)

	for _, schema := range t.all {
		for _, elm := range schema.Elements {
			t.parseElement(elm)
		}

		for _, ct := range schema.ComplexTypes {
			t.parseComplexType(ct)
		}

		for _, st := range schema.SimpleType {
			t.parseSimpleType(st)
		}
	}

	// Return found element name if given type is used only once
	if len(t.foundElementName) > 0 && !t.typeUsageConflict {
		return t.foundElementName
	}

	// Return original type name
	// No element found or conflicting element names found
	return t.typeName
}
