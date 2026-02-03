package hmc

import (
	"cmp"
	"encoding/xml"
	"errors"
	"iter"
	"net/url"
)

type Option struct {
	Label    string `json:"-"`
	Value    string `json:"value"`
	Selected bool   `json:"selected,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func (o Option) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{Name: xml.Name{Local: "c:Option"}}
	if o.Selected {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "selected"}})
	}
	if o.Disabled {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "disabled"}})
	}
	label := cmp.Or(o.Label, o.Value)
	if o.Label != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: o.Value})
	}
	return e.EncodeElement(label, start)
}

type Select struct {
	Multiple bool     `json:"-"`
	Label    string   `json:"-"`
	Name     string   `json:"-"`
	Error    string   `json:"error"`
	Required bool     `json:"-"`
	Options  []Option `json:"options"`
}

var ErrSelectHasNonOption = errors.New("SelectHasNonOption")

// SetValues returns an error if a value is provided that is not listed
// in s.Options. Ignore this error if this Select is meant to allow
// unlisted selections.
func (s *Select) SetValues(values ...string) (err error) {
	for i := range s.Options {
		s.Options[i].Selected = false
	}
	for _, v := range values {
		found := false
		for i, o := range s.Options {
			if o.Value == v && !o.Disabled {
				s.Options[i].Selected = true
				found = true
			}
		}
		if !found {
			err = ErrSelectHasNonOption
			s.Options = append([]Option{{
				Value:    v,
				Selected: true,
			}}, s.Options...)
		}
	}
	return
}

// Values returns a iterator of the values of all selected non-disabled options.
//
// If o.Multiple is false, only the first selected value is returned.
func (s Select) Values() iter.Seq[string] {
	return iter.Seq[string](func(yield func(string) bool) {
		for _, o := range s.Options {
			if !o.Selected || o.Disabled {
				continue
			}
			if !yield(o.Value) || !s.Multiple {
				return
			}
		}
	})
}

// Returns the value of the first selected non-disabled option in s.Options.
func (s Select) Value() string {
	next, stop := iter.Pull(s.Values())
	defer stop()
	val, _ := next()
	return val
}

// ExtractFormValue behaves similarly to [Input.ExtractFormValue]. If s.Multiple is set, all values are taken; if not, the first value is taken.
//
// An error is returned if a value is extracted that is not listed
// in s.Options; but it is safe to ignore this error if unlisted
// selections are allowed. See [Input.SetValues]
func (s *Select) ExtractFormValue(form url.Values) (err error) {
	formValue, ok := form[s.Name]
	if !ok {
		return
	}
	if s.Multiple {
		err = s.SetValues(formValue...)
		delete(form, s.Name)
	} else {
		err = s.SetValues(formValue[0])
		if len(formValue[1:]) > 0 {
			form[s.Name] = formValue[1:]
		} else {
			delete(form, s.Name)
		}
	}
	return
}

func (i Select) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "c:Select"

	if i.Multiple {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "multiple"}})
	}
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "label"}, Value: i.Label})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
	if i.Required {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "required"}, Value: "true"})
	}

	if err := e.EncodeToken(start); err != nil {
		return nil
	}

	for _, o := range i.Options {
		if err := e.EncodeElement(o, start); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}
