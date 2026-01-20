package hyprctl_test

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	"github.com/Teajey/hyprctl"
	"github.com/Teajey/hyprctl/internal/assert"
)

type content struct {
	Username     hyprctl.Input
	Password     hyprctl.Input
	RegisterLink hyprctl.Link
	Submit       hyprctl.Submit
}

func TestSnapshotForm(t *testing.T) {
	form := hyprctl.Form[content]{
		Method: "POST",
		Elements: content{
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
			RegisterLink: hyprctl.Link{
				Name: "Register",
				Href: "/register",
			},
			Submit: hyprctl.Submit{
				Label: "Login",
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
	tm, err := template.ParseFiles("./examples/templates/input.gotmpl")
	assert.FatalErr(t, "parsing template", err)

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
	err = tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotSelect(t *testing.T) {
	tm, err := template.ParseFiles("./examples/templates/input.gotmpl")
	assert.FatalErr(t, "parsing template", err)

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
	err = tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}

func TestSnapshotMultiSelect(t *testing.T) {
	tm, err := template.ParseFiles("./examples/templates/input.gotmpl")
	assert.FatalErr(t, "parsing template", err)

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
	err = tm.ExecuteTemplate(buf, "input", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
}
