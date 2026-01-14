package hyprctl

import (
	"encoding/xml"
)

// Submit represents the control that initiates a Form submission.
//
// A single form may contain many submit options with different names and values.
type Submit struct {
	Label string
	Name  string `json:",omitempty"`
	Value string `json:",omitempty"`
}

func (i Submit) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "Submit"}

	if i.Name != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Name"}, Value: i.Name})
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Value"}, Value: i.Value})
	} else if i.Value != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Value"}, Value: i.Value})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.CharData(i.Label)); err != nil {
		return err
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}
