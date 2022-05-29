package xsd

// RestrictionValue represents a restriction value.
type RestrictionValue struct {
	Doc   string `xml:"annotation>documentation"`
	Value string `xml:"value,attr"`
}
