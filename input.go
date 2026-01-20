package hyprctl

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type Option struct {
	Value    string
	Label    string `json:",omitempty"`
	Selected bool   `json:",omitempty"`
	Disabled bool   `json:",omitempty"`
}

func (o Option) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if o.Selected && o.Value != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Selected"}})
	}
	if o.Disabled {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Disabled"}})
	}
	if o.Label != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Value"}, Value: o.Value})
		if err := e.EncodeToken(start); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.CharData(o.Label)); err != nil {
			return err
		}
		if err := e.EncodeToken(start.End()); err != nil {
			return err
		}
	} else {
		if err := e.EncodeToken(start); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.CharData(o.Value)); err != nil {
			return err
		}
		if err := e.EncodeToken(start.End()); err != nil {
			return err
		}
	}
	return nil
}

// Input describes a piece of data the server needs from the client,
// including validation requirements.
//
// It is analogous to most input elements
// in HTML including <select>, <input> and <textarea>. One Input may
// map to many <input>s that carry the same name attribute.
//
// IMPORTANT: Some fields are mutually-irrelevant; such as Options and MinLength,
// but they are both kept in this struct for simplicity. It is not an error
// to have them both set at the same time, but it is semantically incorrect
// and may cause confusion.
type Input struct {
	Label    string
	Type     string
	Name     string
	Required bool
	Value    string
	// extraValues accommodates inputs such as HTML's checkboxes and multiselect.
	//
	// Use Values() to get all values including Value.
	extraValues []string
	MinLength   uint
	MaxLength   uint
	Step        float32
	Min         string
	Max         string
	Error       string
	Multiple    bool
	Options     []Option
}

func (i *Input) IsSelect() bool {
	return len(i.Options) > 0
}

func (i *Input) IsMultiSelect() bool {
	return i.Multiple && i.IsSelect()
}

// Values returns Value prepended to ExtraValues.
func (i *Input) Values() []string {
	return append([]string{i.Value}, i.extraValues...)
}

// SetValues sets Value and ExtraValues in the same order as [Input.Values].
//
// It also sets p.Options[].Selected when Options are present.
func (p *Input) SetValues(values ...string) {
	length := len(values)
	switch {
	case length > 1:
		p.Value = values[0]
		p.extraValues = values[1:]
	case length > 0:
		p.Value = values[0]
		p.extraValues = nil
	default:
		p.Value = ""
		p.extraValues = nil
	}

	if !p.IsSelect() {
		return
	}

	for _, v := range values {
		found := false
		for i, o := range p.Options {
			if o.Value == v {
				p.Options[i].Selected = true
				found = true
			}
		}
		if !found {
			p.Options = append([]Option{{
				Value:    v,
				Selected: true,
			}}, p.Options...)
		}
	}
}

func (i Input) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	isSelect := i.IsSelect()
	if isSelect {
		if i.IsMultiSelect() {
			start.Name = xml.Name{Local: "MultiSelect"}
		} else {
			start.Name = xml.Name{Local: "Select"}
		}
	} else {
		start.Name = xml.Name{Local: "Input"}
	}

	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Label"}, Value: i.Label})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Name"}, Value: i.Name})
	if !isSelect && i.Type != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Type"}, Value: i.Type})
	}
	if !isSelect {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Value"}, Value: i.Value})
	}
	if i.MinLength > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "MinLength"}, Value: fmt.Sprintf("%d", i.MinLength)})
	}
	if i.MaxLength > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "MaxLength"}, Value: fmt.Sprintf("%d", i.MaxLength)})
	}
	if i.Step > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Step"}, Value: fmt.Sprintf("%f", i.Step)})
	}
	if i.Min != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Min"}, Value: i.Min})
	}
	if i.Max != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Max"}, Value: i.Max})
	}
	if i.Required {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Required"}, Value: "true"})
	}

	if err := e.EncodeToken(start); err != nil {
		return nil
	}

	if i.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "Error"}}
		if err := e.EncodeElement(i.Error, errorStart); err != nil {
			return err
		}

		if len(i.Options) > 0 {
			optionsStart := xml.StartElement{Name: xml.Name{Local: "Options"}}
			if err := e.EncodeToken(optionsStart); err != nil {
				return err
			}
			for _, o := range i.Options {
				optionStart := xml.StartElement{Name: xml.Name{Local: "Option"}}
				if err := e.EncodeElement(o, optionStart); err != nil {
					return err
				}
			}
			if err := e.EncodeToken(optionsStart.End()); err != nil {
				return err
			}
		}
	} else {
		for _, o := range i.Options {
			optionStart := xml.StartElement{Name: xml.Name{Local: "Option"}}
			if err := e.EncodeElement(o, optionStart); err != nil {
				return err
			}
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

// Validate performs some basic checks on the value
// of the input according to its settings.
//
// [Input.Required], [Input.MaxLength], and [Input.MinLength] are checked. Similar to the minimal
// checks that a browser would make for equivalent HTML.
//
// This functionality can be extended with more bespoke validation by
// checking fields and setting the [Input.Error] field accordingly.
func (p *Input) Validate() {
	if p.Required && p.Value == "" {
		p.Error = fmt.Sprintf("%#v is required", p.Name)
	}

	if p.IsSelect() {
		return
	}

	valueLen := len(p.Value)
	if p.MinLength > 0 && int(p.MinLength) > valueLen {
		p.Error = fmt.Sprintf("%#v requires at least %#v characters (currently %#v characters)", p.Name, p.MinLength, valueLen)
	}

	if p.MaxLength > 0 && int(p.MaxLength) < valueLen {
		p.Error = fmt.Sprintf("%#v supports at most %#v characters (currently %#v characters)", p.Name, p.MaxLength, valueLen)
	}
}

// ValueFromUrlValues will search for the Input's value under
// p.Name, setting p.Values.
//
// The found value is deleted from form.
func (p *Input) ValueFromUrlValues(form url.Values) {
	formValue, ok := form[p.Name]
	if ok {
		p.SetValues(formValue...)
		delete(form, p.Name)
	}
}
