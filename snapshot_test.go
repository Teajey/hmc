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
	tm = template.Must(template.New("").ParseFiles("./examples/templates/input.gotmpl"))
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
	FavouriteFood   hyprctl.Input
	Misc            hyprctl.Map
	Login           hyprctl.Link
	Submit          hyprctl.Submit
}

func TestSnapshotForm(t *testing.T) {
	form := myPage{
		Namespace: hyprctl.NewNamespace(),
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
				FavouriteFood: hyprctl.Input{
					Label: "Favourite food",
					Name:  "favFood",
					Options: []hyprctl.Option{
						{Value: "fruit"},
						{Value: "vegetables"},
						{Value: "meat"},
						{Value: "fish"},
						{Value: "bugs"},
					},
				},
				Misc: hyprctl.Map{
					Label: "Any other arbitrary information you wanna provide?",
					Name:  "misc",
				},
				Login: hyprctl.Link{
					Name: "Register",
					Href: "/register",
				},
				Submit: hyprctl.Submit{
					Label: "Login",
				},
			},
		},
	}

	assert.SnapshotXml(t, form)
	assert.SnapshotJson(t, form)
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
	input := hyprctl.Input{
		Label:    "Mug size",
		Type:     "text",
		Name:     "mugs",
		Required: true,
		Options: []hyprctl.Option{
			{Label: "Large", Value: "lg"},
			{Label: "Medium", Value: "md"},
			{Label: "Small", Value: "sm"},
		},
	}
	input.SetValues("Wumbo")
	input.Validate()

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotMultiSelect(t *testing.T) {
	input := hyprctl.Input{
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
	input.Validate()

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotMap(t *testing.T) {
	input := hyprctl.Map{Name: "data"}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValues(form)

	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "only unmatched entries are still in form", 1, len(form))
}

func TestSnapshotBucket(t *testing.T) {
	input := hyprctl.Map{}
	form := url.Values{
		"tree":         {"oak"},
		"data[food]":   {"icecream"},
		"data[drinks]": {"water", "tea"},
	}

	input.ExtractFormValues(form)

	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "all entries are extracted by Map", 0, len(form))
}
