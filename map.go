package hmc

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Map represents an arbitrary set of key-value entries. Where Name == "foo", the form submission may populate Entries with values like this `foo[x]=y&foo[a]=b&foo[stuff]=etc`.
//
// When Name == "", [Map.ExtractFormValue] will extract all values from the given form. Extract specific fields first to prevent them from being captured by the catch-all.
type Map struct {
	Label string `json:"-"`
	Name  string `json:"-"`
	Error string `json:"error"`
	// MaxEntries sets the maximum number of entries allowed when [Map.Validate] is called.
	//
	// Does nothing if set to zero or less.
	MaxEntries int `json:"-"`
	// MaxKeyLength sets the maximum length that any given key may be when [Map.Validate] is called.
	//
	// Does nothing if set to zero or less.
	MaxKeyLength int `json:"-"`
	// MaxValues sets the maximum number of values allowed per entry when [Map.Validate] is called.
	//
	// Does nothing if set to zero or less.
	MaxValues int `json:"-"`
	// MaxValueLength sets the maximum length any given value may be, in any given entry, when [Map.Validate] is called.
	//
	// Does nothing if set to zero or less.
	MaxValueLength int                 `json:"-"`
	Entries        map[string][]string `json:"entries"`
}

// NamedKey returns key as a form name.
// E.g. where m.Name = "foo", m.NamedKey("bar") returns "foo[bar]"
//
// If Name is not set, key is returned unchanged.
// E.g. where m.Name = "", m.NamedKey("bar") returns "bar"
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
	return fmt.Sprintf("contains more than %d entrie(s)", e.Max)
}

type ErrMapMaxLength struct {
	MaxLength int
}

func (e ErrMapMaxLength) Error() string {
	return fmt.Sprintf("contains an entry with more than %d value(s)", e.MaxLength)
}

type ErrMapMaxKeyLength struct {
	MaxKeyLength int
}

func (e ErrMapMaxKeyLength) Error() string {
	return fmt.Sprintf("contains a key longer than %d char(s)", e.MaxKeyLength)
}

type ErrMapMaxValueLength struct {
	Key            string
	MaxValueLength int
}

func (e ErrMapMaxValueLength) Error() string {
	return fmt.Sprintf("key %#v contains a value longer than %d char(s)", e.Key, e.MaxValueLength)
}

// ExtractFormValue takes all entries from form with name x[y], where x = m.Name, and y is an arbitrary key provided by the request.
//
// If m.Name == "", then all entries in form are transferred to m.Entries. Extract specific fields first to prevent them from being captured by the catch-all.
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

// Validate performs some basic checks on m.Entries
// according to given settings.
//
// An error will be returned, and m.Error set, if any of m.MaxEntries, m.MaxKeyLength, m.MaxValueLength, or m.MaxValues are violated.
func (m *Map) Validate() (err error) {
	for k, v := range m.Entries {
		if m.MaxKeyLength > 0 && len(k) > m.MaxKeyLength {
			err = ErrMapMaxKeyLength{m.MaxKeyLength}
		}
		if m.MaxValues > 0 && len(v) > m.MaxValues {
			err = ErrMapMaxLength{m.MaxValues}
		}
		if m.MaxValueLength > 0 {
			for _, val := range v {
				if len(val) > m.MaxValueLength {
					err = ErrMapMaxValueLength{
						Key:            k,
						MaxValueLength: m.MaxValueLength,
					}
				}
			}
		}
	}

	if m.MaxEntries > 0 && len(m.Entries) > m.MaxEntries {
		err = ErrMapMax{m.MaxEntries}
	}

	if err != nil {
		m.Error = err.Error()
	}

	return
}

func (m Map) MarshalXML(e *xml.Encoder, label xml.StartElement) error {
	label.Name.Local = "c:Label"

	if err := e.EncodeToken(label); err != nil {
		return fmt.Errorf("encoding label start: %w", err)
	}
	if err := e.EncodeToken(xml.CharData(m.Label)); err != nil {
		return fmt.Errorf("encoding label text: %w", err)
	}

	inner := xml.StartElement{Name: xml.Name{Local: "c:Map"}}

	inner.Attr = append(inner.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: m.Name})
	if m.MaxEntries > 0 {
		inner.Attr = append(inner.Attr, xml.Attr{Name: xml.Name{Local: "maxentries"}, Value: fmt.Sprintf("%d", m.MaxEntries)})
	}
	if m.MaxKeyLength > 0 {
		inner.Attr = append(inner.Attr, xml.Attr{Name: xml.Name{Local: "maxkeylength"}, Value: fmt.Sprintf("%d", m.MaxKeyLength)})
	}
	if m.MaxValues > 0 {
		inner.Attr = append(inner.Attr, xml.Attr{Name: xml.Name{Local: "maxvalues"}, Value: fmt.Sprintf("%d", m.MaxValues)})
	}
	if m.MaxValueLength > 0 {
		inner.Attr = append(inner.Attr, xml.Attr{Name: xml.Name{Local: "maxvaluelength"}, Value: fmt.Sprintf("%d", m.MaxValueLength)})
	}

	if err := e.EncodeToken(inner); err != nil {
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

	if err := e.EncodeToken(inner.End()); err != nil {
		return fmt.Errorf("encoding map end: %w", err)
	}

	if err := e.EncodeToken(label.End()); err != nil {
		return fmt.Errorf("encoding label end: %w", err)
	}

	return nil
}
