package hmc

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const (
	MaxMapKeyLength = 256
)

// Map represents an arbitrary set of key-value entries. Where Name: "foo", the form submission may provide `foo[x]=y&foo[a]=b&foo[stuff]=etc`.
//
// When name is not set, [Map.ExtractFormValue] will extract all values from the given form.
//
// Key lengths in Entries are limited to [MaxMapKeyLength].
type Map struct {
	Label string `json:"-"`
	Name  string `json:"-"`
	Error string `json:"error"`
	// Max is the maximum amount of entries allowed.
	Max int `json:"-"`
	// MaxLength is the maximum number of values allowed per entry.
	MaxLength int                 `json:"-"`
	Entries   map[string][]string `json:"entries"`
}

func (m Map) NamedKey(key string) string {
	if m.Name == "" {
		return key
	} else {
		return fmt.Sprintf("%s[%s]", m.Name, key)
	}
}

type ErrMapMax struct {
	Max int
}

func (e ErrMapMax) Error() string {
	return fmt.Sprintf("contains more than %d entries", e.Max)
}

type ErrMapMaxLength struct {
	MaxLength int
}

func (e ErrMapMaxLength) Error() string {
	return fmt.Sprintf("contains an entry with more than %d values(s)", e.MaxLength)
}

type ErrMapMaxKeyLength struct {
	MaxLength int
}

func (e ErrMapMaxKeyLength) Error() string {
	return fmt.Sprintf("contains a key longer than %d chars", MaxMapKeyLength)
}

func (m *Map) ExtractFormValue(form url.Values) (err error) {
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
		if len(key) > MaxMapKeyLength {
			err = ErrMapMaxKeyLength{}
		}
		if m.MaxLength > 0 && len(v) > m.MaxLength {
			err = ErrMapMaxLength{m.MaxLength}
		}
		delete(form, k)
		m.Entries[key] = v
	}

	if m.Max > 0 && len(m.Entries) > m.Max {
		err = ErrMapMax{m.Max}
	}

	if err != nil {
		m.Error = err.Error()
	}

	return
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "c:Map"
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "label"}, Value: m.Label})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: m.Name})
	if m.Max > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "max"}, Value: fmt.Sprintf("%d", m.Max)})
	}
	if m.MaxLength > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "maxlength"}, Value: fmt.Sprintf("%d", m.MaxLength)})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if m.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "c:Error"}}
		if err := e.EncodeElement(m.Error, errorStart); err != nil {
			return err
		}
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
