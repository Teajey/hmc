package hmc_test

import (
	"bytes"
	"fmt"
	"html/template"
	"maps"
	"net/url"
	"regexp"
	"slices"
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
	LikesMovies     hmc.Input
	LikesMusic      hmc.Input
	LikesGames      hmc.Input
	TermsAndConds   hmc.Input
	NewsLetter      hmc.Input
	Misc            hmc.Map
	Login           hmc.Link `json:"LoginLink"`
}

func (f *login) controls() []hmc.Control {
	c := []hmc.Control{
		&f.Username,
		&f.Password,
		&f.ConfirmPassword,
		&f.FavouriteFood,
		&f.LikesMovies,
		&f.LikesMusic,
		&f.LikesGames,
		&f.TermsAndConds,
		&f.NewsLetter,
		&f.Misc,
	}
	return c
}

func (f *login) ExtractValues(form url.Values) {
	for _, c := range f.controls() {
		_ = c.ExtractValue(form)
	}
}

func (f *login) Validate() {
	for _, c := range f.controls() {
		_ = c.Validate()
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
				LikesMovies: hmc.Input{
					Label: "Do you like movies?",
					Type:  "checkbox",
					Name:  "likes",
					Value: "movies",
				},
				LikesMusic: hmc.Input{
					Label: "Do you like music?",
					Type:  "checkbox",
					Name:  "likes",
					Value: "music",
				},
				LikesGames: hmc.Input{
					Label: "Do you like games?",
					Type:  "checkbox",
					Name:  "likes",
					Value: "games",
				},
				TermsAndConds: hmc.Input{
					Label: "Do you agree to our terms?",
					Type:  "checkbox",
					Name:  "terms",
					Value: "agreed",
				},
				NewsLetter: hmc.Input{
					Label:   "Do you want our newsletter?",
					Type:    "checkbox",
					Name:    "newsletter",
					Value:   "yes",
					Checked: true,
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
		"likes":            {"music", "movies"},
		"terms":            {"agreed"},
		"misc[iq]":         {"80"},
	}
	page.Form.Elements.ExtractValues(form)
	page.Form.Elements.Validate()

	assert.SnapshotXml(t, page)
	assert.SnapshotJson(t, page)
	unmatched := slices.Collect(maps.Keys(form))
	slices.Sort(unmatched)
	assert.SlicesEq(t, "only unmatched entries remain", []string{"confirm_password", "username"}, unmatched)
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

func TestSnapshotInputDisabled(t *testing.T) {
	input := hmc.Input{
		Label:     "Message",
		Type:      "text",
		Name:      "msg",
		Disabled:  true,
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

func TestSnapshotInputNoLabel(t *testing.T) {
	input := hmc.Input{
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

func TestSnapshotInputPattern(t *testing.T) {
	input := hmc.Input{
		Label:    "Message",
		Type:     "text",
		Name:     "msg",
		Required: true,
		Value:    "Hey...",
		Pattern:  regexp.MustCompile(`\w\.$`),
	}

	_ = input.Validate()

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
	err := input.ExtractValue(form)
	if err != nil {
		input.Error = err.Error()
	}

	buf := bytes.NewBuffer([]byte{})
	err = tm.ExecuteTemplate(buf, "select", input)
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

	_ = input.SetValues("dog", "cat", "mouse")

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

	input.ExtractValue(form)

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

	input.ExtractValue(form)

	buf := bytes.NewBuffer([]byte{})
	err := tm.ExecuteTemplate(buf, "map.gotmpl", input)
	assert.FatalErr(t, "executing template", err)

	assert.Snapshot(t, fmt.Sprintf("%s.snap.html", t.Name()), buf.Bytes())
	assert.SnapshotXml(t, input)
	assert.SnapshotJson(t, input)
	assert.Eq(t, "all entries are extracted by Map", 0, len(form))
}
