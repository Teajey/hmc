package hyprctl_test

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"testing"

	"github.com/Teajey/hyprctl"
	"github.com/Teajey/hyprctl/internal/assert"
)

var tm *template.Template

func TestMain(m *testing.M) {
	tm = template.Must(template.New("").ParseGlob("./examples/templates/*.gotmpl"))
	m.Run()
}

type myPage struct {
	hyprctl.Namespace
	Title string
	Form  hyprctl.Form[login]
}

type login struct {
	Username        hyprctl.Input
	Password        hyprctl.Input
	ConfirmPassword hyprctl.Input
	FavouriteFood   hyprctl.Select
	Misc            hyprctl.Map
	Login           hyprctl.Link
}

func (l *login) ExtractValues(form url.Values) {
	l.Username.ExtractFormValue(form)
	l.Password.ExtractFormValue(form)
	l.ConfirmPassword.ExtractFormValue(form)
	l.FavouriteFood.ExtractFormValues(form)
	l.Misc.ExtractFormValues(form)
}

func (l *login) Validate() {
	l.Username.Validate()
	l.Password.Validate()
	l.ConfirmPassword.Validate()
}

func TestSnapshotForm(t *testing.T) {
	page := myPage{
		Namespace: hyprctl.SetNamespace(),
		Title:     "Login to my thing",
		Form: hyprctl.Form[login]{
			Method: "POST",
			Elements: login{
				Username: hyprctl.Input{
					Label:    "Username",
					Name:     "username",
					Required: true,
				},
				Password: hyprctl.Input{
					Label:    "Password",
					Name:     "password",
					Type:     "password",
					Required: true,
				},
				ConfirmPassword: hyprctl.Input{
					Label:    "Confirm password",
					Name:     "confirmPassword",
					Type:     "password",
					Required: true,
				},
				FavouriteFood: hyprctl.Select{
					Label: "Favourite food",
					Name:  "favFood",
					Options: []hyprctl.Option{
						{Selected: true},
						{Value: "fruit"},
						{Value: "vegetables"},
						{Value: "meat"},
						{Value: "fish"},
						{Label: "Bugs", Value: "bugs"},
					},
					Required: true,
				},
				Misc: hyprctl.Map{
					Label: "Any other arbitrary information you wanna provide?",
					Name:  "misc",
				},
				Login: hyprctl.Link{
					Name: "Register",
					Href: "/register",
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
	link := hyprctl.Link{
		Name: "Register",
		Href: "/register",
	}

	assert.SnapshotXml(t, link)
	assert.SnapshotJson(t, link)
}

func TestSnapshotInput(t *testing.T) {
	input := hyprctl.Input{
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
	input := hyprctl.Select{
		Label:    "Mug size",
		Name:     "mugs",
		Required: true,
		Options: []hyprctl.Option{
			{Label: "Large", Value: "lg"},
			{Label: "Medium", Value: "md"},
			{Label: "Small", Value: "sm"},
		},
	}
	form := url.Values{
		"mugs":  {"Wumbo"},
		"other": {"1"},
	}
	input.ExtractFormValues(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "select", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "only unmatched entries remain", 1, len(form))
}

func TestSnapshotMultiSelect(t *testing.T) {
	input := hyprctl.Select{
		Label:    "Favourite animals",
		Multiple: true,
		Name:     "fav_anim",
		Required: true,
		Options: []hyprctl.Option{
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
	input := hyprctl.Map{Label: "Random data", Name: "data"}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValues(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "map.gotmpl", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "only unmatched entries are still in form", 1, len(form))
}

func TestSnapshotBucket(t *testing.T) {
	input := hyprctl.Map{Label: "Leftover data"}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValues(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "map.gotmpl", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "all entries are extracted by Map", 0, len(form))
}
