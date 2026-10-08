package hmc

import (
	"net/url"
)

// Control represents either an [Input],
// [Select], or [Map]; and their common
// functions.
//
// It's useful for writing batch operations on
// all such controls of a form, e.g.:
//
// func (f *MyForm) Controls() []hmc.Control
//
// hmc.Form[MyForm]
type Control interface {
	// Validate checks the value of the underlying
	// control and sets it's error field accordingly
	Validate() error
	// ExtractValue removes the respective entry from
	// form and adds it to the control's value
	ExtractValue(form url.Values) error
	// AppendValue copies the controls value field
	// to it's respective entry in form.
	AppendValue(form url.Values)
	// SetErrorMsg sets the error field of the underlying control.
	SetErrorMsg(msg string)
	// SetErrorMsg gets the value of the error field of the underlying control.
	ErrorMsg() string
	// FieldName gets the name of the underlying control.
	//
	// e.g. <input name="foo">
	FieldName() string
}

var _ Control = (*Input)(nil)
var _ Control = (*Select)(nil)
var _ Control = (*Map)(nil)

// SetErrorMsg implements [Control.SetErrorMsg]
func (i *Input) SetErrorMsg(msg string) {
	i.Error = msg
}

// SetErrorMsg implements [Control.SetErrorMsg]
func (s *Select) SetErrorMsg(msg string) {
	s.Error = msg
}

// SetErrorMsg implements [Control.SetErrorMsg]
func (m *Map) SetErrorMsg(msg string) {
	m.Error = msg
}

// ErrorMsg implements [Control.ErrorMsg]
func (i *Input) ErrorMsg() string {
	return i.Error
}

// ErrorMsg implements [Control.ErrorMsg]
func (s *Select) ErrorMsg() string {
	return s.Error
}

// ErrorMsg implements [Control.ErrorMsg]
func (m *Map) ErrorMsg() string {
	return m.Error
}

// FieldName implements [Control.FieldName]
func (i *Input) FieldName() string {
	return i.Name
}

// FieldName implements [Control.FieldName]
func (s *Select) FieldName() string {
	return s.Name
}

// FieldName implements [Control.FieldName]
func (m *Map) FieldName() string {
	return m.Name
}
