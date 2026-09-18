package hmc

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"time"
	"unicode/utf8"
)

// Input describes a piece of data the server needs from the client,
// including validation requirements.
//
// It is analogous to HTML's <input> and <textarea>.
type Input struct {
	Label     string
	Type      string
	Name      string
	Value     string
	Checked   bool
	Error     string
	Required  bool
	Disabled  bool
	MinLength uint
	MaxLength uint
	Step      float32
	Min       string
	Max       string
}

func (i Input) MarshalJSON() ([]byte, error) {
	if i.Disabled {
		j := inputJsonDisabled{
			Label:    i.Label,
			Name:     i.Name,
			Disabled: true,
		}
		return json.Marshal(j)
	}
	j := inputJson(i)
	if i.Type == "password" && j.Value != "" {
		j.Value = "********"
	}
	return json.Marshal(j)
}

func (i Input) MarshalXML(e *xml.Encoder, label xml.StartElement) error {
	label.Name.Local = "c:Label"

	if err := e.EncodeToken(label); err != nil {
		return fmt.Errorf("encoding label start: %w", err)
	}
	if err := e.EncodeToken(xml.CharData(i.Label)); err != nil {
		return fmt.Errorf("encoding label text: %w", err)
	}

	input := xml.StartElement{Name: xml.Name{Local: "c:Input"}}

	if i.Disabled {
		input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
		input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "disabled"}, Value: "true"})
	} else {
		if i.Type != "" {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "type"}, Value: i.Type})
		}
		input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: i.Name})
		if i.Type == "password" && i.Value != "" {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: "********"})
		} else {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: i.Value})
		}
		if i.Checked {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "checked"}, Value: "true"})
		}
		if i.MinLength > 0 {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "minlength"}, Value: fmt.Sprintf("%d", i.MinLength)})
		}
		if i.MaxLength > 0 {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "maxlength"}, Value: fmt.Sprintf("%d", i.MaxLength)})
		}
		if i.Step > 0 {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "step"}, Value: fmt.Sprintf("%f", i.Step)})
		}
		if i.Min != "" {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "min"}, Value: i.Min})
		}
		if i.Max != "" {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "max"}, Value: i.Max})
		}
		if i.Required {
			input.Attr = append(input.Attr, xml.Attr{Name: xml.Name{Local: "required"}, Value: "true"})
		}
	}

	if err := e.EncodeToken(input); err != nil {
		return fmt.Errorf("encoding input start: %w", err)
	}

	if i.Error != "" {
		errorStart := xml.StartElement{Name: xml.Name{Local: "c:Error"}}
		if err := e.EncodeElement(i.Error, errorStart); err != nil {
			return fmt.Errorf("encoding error: %w", err)
		}
	}

	if err := e.EncodeToken(input.End()); err != nil {
		return fmt.Errorf("encoding input end: %w", err)
	}

	if err := e.EncodeToken(label.End()); err != nil {
		return fmt.Errorf("encoding label end: %w", err)
	}

	return nil
}

type ErrInputRequired struct{}

func (e ErrInputRequired) Error() string {
	return "value is required"
}

type ErrInputMax struct {
	Max string
}

func (e ErrInputMax) Error() string {
	return fmt.Sprintf("must be at most %s", e.Max)
}

type ErrInputMin struct {
	Min string
}

func (e ErrInputMin) Error() string {
	return fmt.Sprintf("must be at least %s", e.Min)
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
// checks that a browser would make for equivalent HTML. If the type is "checkbox" or "radio" only [Input.Required] is checked.
//
// Because it is complex to implement for whatever many types [Input.Type] might be set to, this function does not validate [Input.Step].
//
// [Input.Type] is treated as a hint to the user, [Input.Value] is not checked against [Input.Type]. Such validation is performed when a corresponding
// Input.ParseValueAs* method is called, although again, [Input.Type] is not checked, and any parse method may be called regardless of the type attribute.
// It is your responsibility to make sure that [Input.Type] is set according to how it is parsed.
//
// This functionality can be extended with more bespoke validation by
// checking fields and setting the [Input.Error] field accordingly.
func (i *Input) Validate() error {
	if i.Type == "checkbox" || i.Type == "radio" {
		if i.Required && !i.Checked {
			return ErrInputRequired{}
		}
		return nil
	}

	if i.Value == "" {
		if i.Required {
			return ErrInputRequired{}
		}
		return nil
	}

	if i.Max != "" && i.cmpLess(i.Max, i.Value) {
		return ErrInputMax{i.Max}
	}
	if i.Min != "" && i.cmpLess(i.Value, i.Min) {
		return ErrInputMin{i.Min}
	}

	valueLen := utf8.RuneCountInString(i.Value)
	if i.MaxLength > 0 && valueLen > int(i.MaxLength) {
		return ErrInputMaxLength{i.MaxLength}
	}
	if i.MinLength > 0 && valueLen < int(i.MinLength) {
		return ErrInputMinLength{i.MinLength}
	}

	return nil
}

// popFirst removes and returns the first value stored at key.
func popFirst(form url.Values, key string) (string, bool) {
	vals, ok := form[key]
	if !ok || len(vals) == 0 {
		return "", false
	}
	setOrDelete(form, key, vals[1:])
	return vals[0], true
}

// popValue removes the first occurrence of val stored at key.
func popValue(form url.Values, key, val string) bool {
	vals, ok := form[key]
	if !ok {
		return false
	}
	idx := slices.Index(vals, val)
	if idx == -1 {
		return false
	}
	remaining := slices.Concat(vals[:idx:idx], vals[idx+1:])
	setOrDelete(form, key, remaining)
	return true
}

func setOrDelete(form url.Values, key string, vals []string) {
	if len(vals) == 0 {
		delete(form, key)
		return
	}
	form[key] = vals
}

// ExtractFormValue sets i.Value to the first value found at form[i.Name],
// removing it from form. For checkboxes and radios it instead sets i.Checked
// and removes the matching value.
func (i *Input) ExtractFormValue(form url.Values) {
	if i.Disabled {
		return
	}
	switch i.Type {
	case "checkbox", "radio":
		i.Checked = popValue(form, i.Name, i.Value)
	default:
		if v, ok := popFirst(form, i.Name); ok {
			i.Value = v
		}
	}
}

type ErrInputValueAsTime struct {
	Err error
}

func (e ErrInputValueAsTime) Error() string {
	return "not a valid time"
}

func (e ErrInputValueAsTime) Unwrap() error {
	return e.Err
}

// ParseValueAsTime parses i.Value as `type="time"`. An ISO 8601 time.
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
	return
}

type ErrInputValueAsDate struct {
	Err error
}

func (e ErrInputValueAsDate) Error() string {
	return "not a valid date"
}

func (e ErrInputValueAsDate) Unwrap() error {
	return e.Err
}

// ParseValueAsDate parses i.Value as `type="date"`. An ISO 8601 date.
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
	return
}

type ErrInputValueAsDatetime struct {
	Err error
}

func (e ErrInputValueAsDatetime) Error() string {
	return "not a valid datetime"
}

func (e ErrInputValueAsDatetime) Unwrap() error {
	return e.Err
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
	return
}

type ErrInputValueAsDatetimeLocal struct {
	Err error
}

func (e ErrInputValueAsDatetimeLocal) Error() string {
	return "not a valid datetime-local"
}

func (e ErrInputValueAsDatetimeLocal) Unwrap() error {
	return e.Err
}

// ParseValueAsDatetimeLocal parses i.Value as `type="datetime-local"`. An ISO 8601 datetime without a timezone.
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
	return
}
