package atompub

import "encoding/xml"

// ServiceDocument represents the AtomPub service document (RFC 5023 8.  Service Documents).
type ServiceDocument struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2007/app service"`
	Workspaces []Workspace `xml:"workspace"`
}

// Workspace represents a workspace in a service document
type Workspace struct {
	Title       Text         `xml:"http://www.w3.org/2005/Atom title"`
	Collections []Collection `xml:"collection"`
}

// Collection represents a collection in a workspace
type Collection struct {
	Href       string      `xml:"href,attr,omitempty"`
	Title      Text        `xml:"http://www.w3.org/2005/Atom title"`
	Accept     []string    `xml:"accept,omitempty"`
	Categories *Categories `xml:"categories,omitempty"`
}

// Categories represents the app:categories element
type Categories struct {
	Fixed  string `xml:"fixed,attr,omitempty"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Href   string `xml:"href,attr,omitempty"`
}
