package hyprctl

import (
	"cmp"
	"encoding/xml"
	"iter"
	"net/url"
)

type Option struct {
	Label    string `json:",omitempty"`
	Value    string
	Selected bool `json:",omitempty"`
	Disabled bool `json:",omitempty"`
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
	Multiple bool     `json:"multiple,omitempty"`
	Label    string   `json:"label"`
	Name     string   `json:"name"`
	Required bool     `json:"required,omitempty"`
	Options  []Option `json:"options"`
	Error    string   `json:"error,omitempty"`
}

func (s *Select) SetValues(values ...string) {
	for i := range s.Options {
		s.Options[i].Selected = false
	}
	for _, v := range values {
		found := false
		for i, o := range s.Options {
			if o.Value == v {
				s.Options[i].Selected = true
				found = true
			}
		}
		if !found {
			s.Options = append([]Option{{
				Value:    v,
				Selected: true,
			}}, s.Options...)
		}
	}
}

func (s *Select) Values() iter.Seq[string] {
	return iter.Seq[string](func(yield func(string) bool) {
		for _, o := range s.Options {
			if o.Selected {
				if !yield(o.Value) {
					return
				}
			}
		}
	})
}

func (s *Select) Value() string {
	next, stop := iter.Pull(s.Values())
	defer stop()
	val, _ := next()
	return val
}

func (s *Select) ExtractFormValues(form url.Values) {
	formValue, ok := form[s.Name]
	if !ok {
		return
	}
	if s.Multiple {
		s.SetValues(formValue...)
		delete(form, s.Name)
	} else {
		s.SetValues(formValue[0])
		if len(formValue[1:]) > 0 {
			form[s.Name] = formValue[1:]
		} else {
			delete(form, s.Name)
		}
	}
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
