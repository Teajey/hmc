package hyprctl

import "encoding/xml"

// Link represents a state transition that requires no input—a simple
// navigation or action trigger.
type Link struct {
	Href string
	Name string
}

func (i Link) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "A"}

	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Href"}, Value: i.Href})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.CharData(i.Name)); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
