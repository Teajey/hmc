package hmc_test

import (
	"iter"
	"maps"
	"net/url"
	"slices"
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

type roundTripForm struct {
	Text   hmc.Input
	Box    hmc.Input
	RadioA hmc.Input
	RadioB hmc.Input
	Multi  hmc.Select
	Rest   hmc.Map
}

func newRoundTripForm() *roundTripForm {
	return &roundTripForm{
		Text:   hmc.Input{Name: "text"},
		Box:    hmc.Input{Name: "box", Type: "checkbox", Value: "on"},
		RadioA: hmc.Input{Name: "radio", Type: "radio", Value: "a"},
		RadioB: hmc.Input{Name: "radio", Type: "radio", Value: "b"},
		Multi:  hmc.Select{Name: "multi", Multiple: true, Options: []hmc.Option{{Value: "x"}, {Value: "y"}}},
		Rest:   hmc.Map{},
	}
}

func (f *roundTripForm) Controls() iter.Seq[hmc.Control] {
	// The catch-all map must come last or it swallows everything.
	return slices.Values([]hmc.Control{&f.Text, &f.Box, &f.RadioA, &f.RadioB, &f.Multi, &f.Rest})
}

// A form's appended values, resubmitted to a fresh form, reproduce themselves.
func TestFormRoundTrip(t *testing.T) {
	in := url.Values{
		"text":  {"hello"},
		"box":   {"on"},
		"radio": {"b"},
		"multi": {"x", "y"},
		"other": {"1"},
	}

	first := newRoundTripForm()
	form := maps.Clone(in)
	for c := range first.Controls() {
		_ = c.ExtractValue(form)
	}
	assert.Eq(t, "first extraction leftover", 0, len(form))

	out := url.Values{}
	for c := range first.Controls() {
		c.AppendValue(out)
	}
	assertValuesEq(t, "first append", in, out)

	second := newRoundTripForm()
	for c := range second.Controls() {
		_ = c.ExtractValue(out)
	}
	assert.Eq(t, "second extraction leftover", 0, len(out))
	assert.True(t, "radio b checked", second.RadioB.Checked && !second.RadioA.Checked)
}
