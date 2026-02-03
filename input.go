package hmc

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
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
	Error     string
	Required  bool
	Value     string
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
		return err
	}

	if i.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "c:Error"}}
		if err := e.EncodeElement(i.Error, errorStart); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type ErrInputRequired struct{}

func (e ErrInputRequired) Error() string {
	return "value is required"
}

type ErrInputMax struct {
	Max string
}

func (e ErrInputMax) Error() string {
	return fmt.Sprintf("must be less than or equal to %s", e.Max)
}

type ErrInputMin struct {
	Min string
}

func (e ErrInputMin) Error() string {
	return fmt.Sprintf("must be greater than %s", e.Min)
}

type ErrInputMaxLength struct {
	MaxLength uint
}

func (e ErrInputMaxLength) Error() string {
	return fmt.Sprintf("must be at most %d char(s)", e.MaxLength)
}

type ErrInputMinLength struct {
	MinLength uint
}

func (e ErrInputMinLength) Error() string {
	return fmt.Sprintf("must be at least %d char(s)", e.MinLength)
}

func (p *Input) cmpLess(x, y string) bool {
	if p.Type == "number" || p.Type == "range" {
		px, xErr := strconv.ParseFloat(x, 64)
		py, yErr := strconv.ParseFloat(y, 64)
		if xErr == nil && yErr == nil {
			return px < py
		}
	}

	return cmp.Less(x, y)
}

// Validate performs some basic checks on the value
// of the input according to its settings.
//
// [Input.Required], [Input.Max], [Input.Min], [Input.MaxLength], and [Input.MinLength] are checked, in that order. Similar to the minimal
// checks that a browser would make for equivalent HTML.
//
// Because it is complex to implement for whatever many types [Input.Type] might be set to, this function does not validate [Input.Step].
//
// [Input.Type] is treated as a hint to the user, [Input.Value] is not checked against [Input.Type]. Such validation is performed when a corresponding
// Input.ParseValueAs* method is called, although again, [Input.Type] is not checked, and any parse method may be called regardless of the type attribute.
// It is your responsibility to make sure that [Input.Type] is set according to how it is parsed.
//
// This functionality can be extended with more bespoke validation by
// checking fields and setting the [Input.Error] field accordingly.
func (p *Input) Validate() (err error) {
	if p.Required && p.Value == "" {
		err = ErrInputRequired{}
	}

	if p.Value != "" {
		if p.Max != "" && p.cmpLess(p.Max, p.Value) {
			err = ErrInputMax{p.Max}
		}

		if p.Min != "" && p.cmpLess(p.Value, p.Min) {
			err = ErrInputMin{p.Min}
		}
	}

	valueLen := len(p.Value)
	if p.MaxLength > 0 && int(p.MaxLength) < valueLen {
		err = ErrInputMaxLength{p.MaxLength}
	}

	if p.MinLength > 0 && valueLen < int(p.MinLength) {
		err = ErrInputMinLength{p.MinLength}
	}

	if err != nil {
		p.Error = err.Error()
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

type ErrInputValueAsTime struct {
	Err error
}

func (e ErrInputValueAsTime) Error() string {
	return "not a valid time"
}

// ParseValueAsDatetime parses i.Value as `type="time"`. An ISO 8601 time.
//
// [Input.Type] is not checked here.
func (i *Input) ParseValueAsTime() (t time.Time, err error) {
	if i.Value == "" {
		return
	}
	t, err = time.Parse("15:04:05.999999999", i.Value)
	if err == nil {
		return
	}
	err = ErrInputValueAsTime{
		err,
	}
	i.Error = err.Error()
	return
}

type ErrInputValueAsDate struct {
	Err error
}

func (e ErrInputValueAsDate) Error() string {
	return "not a valid date"
}

// ParseValueAsDatetime parses i.Value as `type="date"`. An ISO 8601 date.
//
// [Input.Type] is not checked here.
func (i *Input) ParseValueAsDate() (t time.Time, err error) {
	if i.Value == "" {
		return
	}
	t, err = time.Parse(time.DateOnly, i.Value)
	if err == nil {
		return
	}
	err = ErrInputValueAsDate{
		err,
	}
	i.Error = err.Error()
	return
}

type ErrInputValueAsDatetime struct {
	Err error
}

func (e ErrInputValueAsDatetime) Error() string {
	return "not a valid datetime"
}

// ParseValueAsDatetime parses i.Value as `type="datetime"`. An ISO 8601 datetime that expects a timezone.
//
// WARNING: This is not widely supported by browsers.
//
// [Input.Type] is not checked here.
func (i *Input) ParseValueAsDatetime() (t time.Time, err error) {
	if i.Value == "" {
		return
	}
	t, err = time.Parse(time.RFC3339Nano, i.Value)
	if err == nil {
		return
	}
	err = ErrInputValueAsDatetime{
		err,
	}
	i.Error = err.Error()
	return
}

type ErrInputValueAsDatetimeLocal struct {
	Err error
}

func (e ErrInputValueAsDatetimeLocal) Error() string {
	return "not a valid datetime-local"
}

// ParseValueAsDatetime parses i.Value as `type="datetime-local"`. An ISO 8601 datetime without a timezone.
//
// [Input.Type] is not checked here.
func (i *Input) ParseValueAsDatetimeLocal() (t time.Time, err error) {
	if i.Value == "" {
		return
	}
	t, err = time.Parse("2006-01-02T15:04:05.999999999", i.Value)
	if err == nil {
		return
	}
	err = ErrInputValueAsDatetimeLocal{
		err,
	}
	i.Error = err.Error()
	return
}
