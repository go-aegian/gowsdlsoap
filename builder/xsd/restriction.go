package xsd

// Restriction defines restrictions on a simpleType, simpleContent, or complexContent definition.
type Restriction struct {
	Base         string             `xml:"base,attr"`
	Enumeration  []RestrictionValue `xml:"enumeration"`
	Pattern      RestrictionValue   `xml:"pattern"`
	MinInclusive RestrictionValue   `xml:"minInclusive"`
	MaxInclusive RestrictionValue   `xml:"maxInclusive"`
	WhiteSpace   RestrictionValue   `xml:"whitespace"`
	Length       RestrictionValue   `xml:"length"`
	MinLength    RestrictionValue   `xml:"minLength"`
	MaxLength    RestrictionValue   `xml:"maxLength"`
}
