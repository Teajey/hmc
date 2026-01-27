package hyprctl

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"net/url"
)

// Input describes a piece of data the server needs from the client,
// including validation requirements.
//
// It is analogous to HTML's <input> and <textarea>.
//
// IMPORTANT: Some fields are mutually-irrelevant; such as Options and MinLength,
// but they are both kept in this struct for simplicity. It is not an error
// to have them both set at the same time, but it is semantically incorrect
// and may cause confusion.
type Input struct {
	Label     string
	Type      string
	Name      string
	Required  bool
	Value     string
	Error     string
	MinLength uint
	MaxLength uint
	Step      float32
	Min       string
	Max       string
}

func (i Input) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "c:Input"

	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "label"}, Value: i.Label})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
	if i.Type != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "type"}, Value: i.Type})
	}
	if i.Type == "password" && i.Value != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: "********"})
	} else {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: i.Value})
	}
	if i.MinLength > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "minlength"}, Value: fmt.Sprintf("%d", i.MinLength)})
	}
	if i.MaxLength > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "maxlength"}, Value: fmt.Sprintf("%d", i.MaxLength)})
	}
	if i.Step > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "step"}, Value: fmt.Sprintf("%f", i.Step)})
	}
	if i.Min != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "min"}, Value: i.Min})
	}
	if i.Max != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "max"}, Value: i.Max})
	}
	if i.Required {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "required"}, Value: "true"})
	}

	if err := e.EncodeToken(start); err != nil {
		return nil
	}

	if i.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "c:Error"}}
		if err := e.EncodeElement(i.Error, errorStart); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// Validate performs some basic checks on the value
// of the input according to its settings.
//
// [Input.Required], [Input.Max], [Input.Min], [Input.MaxLength], and [Input.MinLength] are checked, in that order. Similar to the minimal
// checks that a browser would make for equivalent HTML.
//
// This functionality can be extended with more bespoke validation by
// checking fields and setting the [Input.Error] field accordingly.
func (p *Input) Validate() {
	if p.Required && p.Value == "" {
		p.Error = fmt.Sprintf("%#v is required", p.Name)
	}

	if p.Max != "" && cmp.Less(p.Max, p.Value) {
		p.Error = fmt.Sprintf("%#v must be less than %#v", p.Value, p.Max)
	}

	if p.Min != "" && cmp.Less(p.Value, p.Min) {
		p.Error = fmt.Sprintf("%#v must be greater than %#v", p.Value, p.Max)
	}

	valueLen := len(p.Value)
	if p.MinLength > 0 && int(p.MinLength) > valueLen {
		p.Error = fmt.Sprintf("%#v requires at least %#v characters (currently %#v characters)", p.Name, p.MinLength, valueLen)
	}

	if p.MaxLength > 0 && int(p.MaxLength) < valueLen {
		p.Error = fmt.Sprintf("%#v supports at most %#v characters (currently %#v characters)", p.Name, p.MaxLength, valueLen)
	}
}

// ValueFromUrlValues will searching for the Input's value under
// p.Name, setting p.Value.
//
// The found value is deleted from form.
func (i *Input) ExtractFormValue(form url.Values) {
	formValue, ok := form[i.Name]
	if ok {
		i.Value = formValue[0]
		if len(formValue[1:]) > 0 {
			form[i.Name] = formValue[1:]
		} else {
			delete(form, i.Name)
		}
	}
}
