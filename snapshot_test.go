package hmc_test

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

var tm *template.Template

func TestMain(m *testing.M) {
	tm = template.Must(template.New("").ParseGlob("./examples/templates/*.gotmpl"))
	m.Run()
}

type myPage struct {
	hmc.Namespace
	Title string
	Form  hmc.Form[login]
}

type login struct {
	Username        hmc.Input
	Password        hmc.Input
	ConfirmPassword hmc.Input
	FavouriteFood   hmc.Select
	Misc            hmc.Map
	Login           hmc.Link `json:"LoginLink"`
}

func (l *login) ExtractValues(form url.Values) {
	l.Username.ExtractFormValue(form)
	l.Password.ExtractFormValue(form)
	l.ConfirmPassword.ExtractFormValue(form)
	_ = l.FavouriteFood.ExtractFormValue(form)
	l.Misc.ExtractFormValue(form)
}

func (l *login) Validate() {
	err := l.Username.Validate()
	if err != nil {
		l.Username.Error = err.Error()
	}
	err = l.Password.Validate()
	if err != nil {
		l.Password.Error = err.Error()
	}
	err = l.ConfirmPassword.Validate()
	if err != nil {
		l.ConfirmPassword.Error = err.Error()
	}
}

func TestSnapshotForm(t *testing.T) {
	page := myPage{
		Namespace: hmc.NS(),
		Title:     "Login to my thing",
		Form: hmc.Form[login]{
			Method: "POST",
			Elements: login{
				Username: hmc.Input{
					Label:    "Username",
					Name:     "username",
					Required: true,
				},
				Password: hmc.Input{
					Label:    "Password",
					Name:     "password",
					Type:     "password",
					Required: true,
				},
				ConfirmPassword: hmc.Input{
					Label:    "Confirm password",
					Name:     "confirmPassword",
					Type:     "password",
					Required: true,
				},
				FavouriteFood: hmc.Select{
					Label: "Favourite food",
					Name:  "favFood",
					Options: []hmc.Option{
						{Selected: true},
						{Value: "fruit"},
						{Value: "vegetables"},
						{Value: "meat"},
						{Value: "fish"},
						{Label: "Bugs", Value: "bugs"},
					},
					Required: true,
				},
				Misc: hmc.Map{
					Label: "Any other arbitrary information you wanna provide?",
					Name:  "misc",
				},
				Login: hmc.Link{
					Label: "Register",
					Href:  "/register",
				},
			},
		},
	}

	form := url.Values{
		"username":         {"john", "blane"},
		"password":         {"123456"},
		"confirm_password": {"123456"},
		"favFood":          {"bugs"},
		"misc[iq]":         {"80"},
	}
	page.Form.Elements.ExtractValues(form)
	page.Form.Elements.Validate()

	assert.SnapshotXml(t, page)
	assert.SnapshotJson(t, page)
	assert.Eq(t, "only unmatched entries remain", 2, len(form))
}

func TestSnapshotLink(t *testing.T) {
	link := hmc.Link{
		Label: "Register",
		Href:  "/register",
	}

	assert.SnapshotXml(t, link)
	assert.SnapshotJson(t, link)
}

func TestSnapshotInput(t *testing.T) {
	input := hmc.Input{
		Label:     "Message",
		Type:      "text",
		Name:      "msg",
		Required:  true,
		Value:     "Hey...",
		MinLength: 3,
		Error:     "This is a bad message",
	}

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotSelect(t *testing.T) {
	input := hmc.Select{
		Label:    "Mug size",
		Name:     "mugs",
		Required: true,
		Options: []hmc.Option{
			{Label: "Large", Value: "lg"},
			{Label: "Medium", Value: "md"},
			{Label: "Small", Value: "sm"},
		},
	}
	form := url.Values{
		"mugs":  {"Wumbo"},
		"other": {"1"},
	}
	input.ExtractFormValue(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "select", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "only unmatched entries remain", 1, len(form))
}

func TestSnapshotMultiSelect(t *testing.T) {
	input := hmc.Select{
		Label:    "Favourite animals",
		Multiple: true,
		Name:     "fav_anim",
		Required: true,
		Options: []hmc.Option{
			{Label: "Dog", Value: "dog"},
			{Label: "Cat", Value: "cat"},
			{Label: "Guinea pig", Value: "gn_pig"},
		},
	}

	input.SetValues("dog", "cat", "mouse")

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "select", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotMap(t *testing.T) {
	input := hmc.Map{Label: "Random data", Name: "data"}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValue(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "map.gotmpl", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "only unmatched entries are still in form", 1, len(form))
}

func TestSnapshotBucket(t *testing.T) {
	input := hmc.Map{Label: "Leftover data"}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValue(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "map.gotmpl", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "all entries are extracted by Map", 0, len(form))
}
