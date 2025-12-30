package atompub

import "encoding/xml"

// ServiceDocument represents the AtomPub service document (RFC 5023 ServiceDocument Document).
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
	Href   string   `xml:"href,attr,omitempty"`
	Title  Text     `xml:"http://www.w3.org/2005/Atom title"`
	Accept []string `xml:"accept,omitempty"`
}
