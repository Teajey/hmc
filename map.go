package hyprctl

import (
	"encoding/xml"
	"net/url"
	"strings"
)

type Map struct {
	Name    string
	Entries map[string][]string
}

func (m *Map) ExtractFormValues(form url.Values) {
	if m.Entries == nil {
		m.Entries = make(map[string][]string, len(form))
	}
	for k, v := range form {
		if m.Name != "" && !strings.HasPrefix(k, m.Name+"[") && !strings.HasSuffix(k, "]") {
			continue
		}
		delete(form, k)
		m.Entries[k] = v
	}
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Map"
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Name"}, Value: m.Name})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, values := range m.Entries {
		for _, v := range values {
			keyElem := xml.StartElement{Name: xml.Name{Local: "Input"}}
			keyElem.Attr = append(keyElem.Attr, xml.Attr{Name: xml.Name{Local: "Name"}, Value: k})
			keyElem.Attr = append(keyElem.Attr, xml.Attr{Name: xml.Name{Local: "Value"}, Value: v})
			if err := e.EncodeElement("", keyElem); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(start.End())
}
