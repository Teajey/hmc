package hmc

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type Map struct {
	Label   string              `json:"-"`
	Name    string              `json:"-"`
	Error   string              `json:"error"`
	Entries map[string][]string `json:"entries"`
}

func (m Map) NamedKey(key string) string {
	if m.Name == "" {
		return key
	} else {
		return fmt.Sprintf("%s[%s]", m.Name, key)
	}
}

func (m *Map) ExtractFormValue(form url.Values) {
	if m.Entries == nil {
		m.Entries = make(map[string][]string, len(form))
	}
	for k, v := range form {
		var key string
		if m.Name != "" {
			after, ok := strings.CutPrefix(k, m.Name+"[")
			if !ok {
				continue
			}
			key, ok = strings.CutSuffix(after, "]")
			if !ok {
				continue
			}
		} else {
			key = k
		}
		delete(form, k)
		m.Entries[key] = v
	}
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "c:Map"
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "label"}, Value: m.Label})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: m.Name})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	keys := make([]string, 0, len(m.Entries))
	for k := range m.Entries {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		for _, v := range m.Entries[k] {
			keyElem := xml.StartElement{Name: xml.Name{Local: "c:Input"}}
			keyElem.Attr = append(keyElem.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: m.NamedKey(k)})
			keyElem.Attr = append(keyElem.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: v})
			if err := e.EncodeElement("", keyElem); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(start.End())
}
