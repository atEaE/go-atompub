package atompub

import (
	"encoding/xml"
	"time"
)

// Entry represents an Atom Entry (RFC 4287)
type Entry struct {
	XMLName    xml.Name   `xml:"http://www.w3.org/2005/Atom entry"`
	ID         string     `xml:"id"`
	Title      Text       `xml:"title"`
	Published  *time.Time `xml:"published,omitempty"`
	Authors    []Person   `xml:"author,omitempty"`
	Content    *Content   `xml:"content,omitempty"`
	Links      []Link     `xml:"link,omitempty"`
	Categories []Category `xml:"category,omitempty"`
	Control    *Control   `xml:"http://www.w3.org/2007/app control,omitempty"`
}

// Category represents a category element
type Category struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr,omitempty"`
}

// Control represents app:control element (RFC 5023)
type Control struct {
	Draft string `xml:"http://www.w3.org/2007/app draft,omitempty"`
}

// IsDraft returns true if the entry is marked as draft
func (c *Control) IsDraft() bool {
	return c != nil && c.Draft == "yes"
}

// Person represents an author or contributor
type Person struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

// Content represents the content element
type Content struct {
	Type  string `xml:"type,attr,omitempty"`
	Src   string `xml:"src,attr,omitempty"`
	Value string `xml:",chardata"`
}

// Link represents an Atom link element
type Link struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   int64  `xml:"length,attr,omitempty"`
}

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
