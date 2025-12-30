package atompub

// Text represents a text construct in Atom
type Text struct {
	Type  TextType `xml:"type,attr,omitempty"`
	Value string   `xml:",chardata"`
}

type TextType string

const (
	TextTypeText TextType = "text"
	TextTypeHTML TextType = "html"
)
