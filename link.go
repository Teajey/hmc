package hmc

import "encoding/xml"

// Link represents a state transition that requires no input—a simple
// navigation or action trigger.
type Link struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

func (i Link) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "c:Link"}

	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "href"}, Value: i.Href})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.CharData(i.Label)); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
