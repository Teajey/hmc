package hmc

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"time"
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

func (i Input) MarshalJSON() ([]byte, error) {
	j := inputJson(i)
	if j.Type == "password" && j.Value != "" {
		j.Value = "********"
	}
	return json.Marshal(j)
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

var ErrInputRequired = errors.New("value is required")
var ErrInputMax = errors.New("value is too great")
var ErrInputMin = errors.New("value is too little")
var ErrInputMaxLength = errors.New("value is too long")
var ErrInputMinLength = errors.New("value is too short")

// Validate performs some basic checks on the value
// of the input according to its settings.
//
// [Input.Required], [Input.Max], [Input.Min], [Input.MaxLength], and [Input.MinLength] are checked, in that order. Similar to the minimal
// checks that a browser would make for equivalent HTML.
//
// This functionality can be extended with more bespoke validation by
// checking fields and setting the [Input.Error] field accordingly.
func (p *Input) Validate() (err error) {
	if p.Required && p.Value == "" {
		err = ErrInputRequired
	}

	if p.Value != "" {
		if p.Max != "" && cmp.Less(p.Max, p.Value) {
			err = ErrInputMax
		}

		if p.Min != "" && cmp.Less(p.Value, p.Min) {
			err = ErrInputMin
		}
	}

	valueLen := len(p.Value)
	if p.MinLength > 0 && int(p.MinLength) > valueLen {
		err = ErrInputMinLength
	}

	if p.MaxLength > 0 && int(p.MaxLength) < valueLen {
		err = ErrInputMaxLength
	}

	return
}

// ExtractFormValue sets i.Value to the first value found at form[i.Name].
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

func (i *Input) ParseValueAsTime() (t time.Time, err error) {
	if i.Value != "" {
		t, err = time.Parse("15:04:05.999999999", i.Value)
	}
	return
}

func (i *Input) ParseValueAsDate() (t time.Time, err error) {
	if i.Value != "" {
		t, err = time.Parse(time.DateOnly, i.Value)
	}
	return
}

func (i *Input) ParseValueAsDatetime() (t time.Time, err error) {
	if i.Value != "" {
		t, err = time.Parse(time.RFC3339Nano, i.Value)
	}
	return
}

func (i *Input) ParseValueAsDatetimeLocal() (t time.Time, err error) {
	if i.Value != "" {
		t, err = time.Parse("2006-01-02T15:04:05.999999999", i.Value)
	}
	return
}
