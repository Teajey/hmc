package hmc

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"iter"
	"net/url"
)

type Option struct {
	Label    string `json:"label,omitempty"`
	Value    string `json:"value"`
	Selected bool   `json:"selected,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func (o Option) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{Name: xml.Name{Local: "c:Option"}}
	if o.Disabled {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "disabled"}})
	}
	label := cmp.Or(o.Label, o.Value)
	if o.Label != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: o.Value})
	}
	if o.Selected {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "selected"}})
	}

	if err := e.EncodeToken(start); err != nil {
		return nil
	}

	if err := e.EncodeToken(xml.CharData(label)); err != nil {
		return err
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

type Select struct {
	Label    string   `json:"label"`
	Multiple bool     `json:"multiple,omitempty"`
	Name     string   `json:"name"`
	Error    string   `json:"error,omitempty"`
	Required bool     `json:"required,omitempty"`
	Options  []Option `json:"options"`
	Disabled bool     `json:"disabled,omitempty"`
}

type ErrSelectHasNonOption struct{}

func (ErrSelectHasNonOption) Error() string {
	return "Unlisted selection"
}

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
			err = ErrSelectHasNonOption{}
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
// See also Value() for easily getting just the first selected value.
func (s Select) Values() iter.Seq[string] {
	return iter.Seq[string](func(yield func(string) bool) {
		for _, o := range s.Options {
			if !o.Selected || o.Disabled || o.Value == "" {
				continue
			}
			if !yield(o.Value) {
				return
			}
		}
	})
}

// Returns the value of the first selected non-disabled option in s.Options.
//
// See also Values() for getting all selected values when s.Multiple == true.
func (s Select) Value() string {
	next, stop := iter.Pull(s.Values())
	defer stop()
	val, _ := next()
	return val
}

// Validate performs some basic checks on p according to its settings.
//
// It currently only checks p.Required.
//
// An error will be returned if p.Value() == "" when p.Required == true
func (p *Select) Validate() (err error) {
	if p.Required && p.Value() == "" {
		err = ErrInputRequired{}
	}

	return
}

// ExtractFormValue behaves similarly to [Input.ExtractFormValue]. If s.Multiple is set, all values are taken; if not, the first value is taken.
//
// An error is returned if a value is extracted that is not listed
// in s.Options; but it is safe to ignore this error if unlisted
// selections are allowed. See [Select.SetValues]
func (s *Select) ExtractFormValue(form url.Values) (err error) {
	if s.Disabled {
		return
	}
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

func (i Select) MarshalXML(e *xml.Encoder, label xml.StartElement) error {
	label.Name.Local = "c:Label"

	if err := e.EncodeToken(label); err != nil {
		return fmt.Errorf("encoding label start: %w", err)
	}
	if err := e.EncodeToken(xml.CharData(i.Label)); err != nil {
		return fmt.Errorf("encoding label text: %w", err)
	}

	sel := xml.StartElement{Name: xml.Name{Local: "c:Select"}}

	if i.Disabled {
		sel.Attr = append(sel.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
	} else {
		if i.Multiple {
			sel.Attr = append(sel.Attr, xml.Attr{Name: xml.Name{Local: "multiple"}, Value: "true"})
		}
		sel.Attr = append(sel.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
		if i.Required {
			sel.Attr = append(sel.Attr, xml.Attr{Name: xml.Name{Local: "required"}, Value: "true"})
		}
	}

	if err := e.EncodeToken(sel); err != nil {
		return nil
	}

	if i.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "c:Error"}}
		if err := e.EncodeElement(i.Error, errorStart); err != nil {
			return fmt.Errorf("encoding error: %w", err)
		}
	}

	for _, o := range i.Options {
		if err := e.EncodeElement(o, sel); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(sel.End()); err != nil {
		return fmt.Errorf("encoding select end: %w", err)
	}

	if err := e.EncodeToken(label.End()); err != nil {
		return fmt.Errorf("encoding label end: %w", err)
	}

	return nil
}
