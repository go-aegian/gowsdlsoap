package xsd

// Group element is used to define a group of elements to be used in complex type definitions.
type Group struct {
	Name     string    `xml:"name,attr"`
	Ref      string    `xml:"ref,attr"`
	Sequence []Element `xml:"sequence>element"`
	Choice   []Element `xml:"choice>element"`
	All      []Element `xml:"all>element"`
}
