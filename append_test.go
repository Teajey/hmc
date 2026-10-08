package hmc_test

import (
	"maps"
	"net/url"
	"slices"
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

// Extraction consumes form, so leftover is whatever the control didn't claim.
func roundTrip(c hmc.Control, in url.Values) (out, leftover url.Values) {
	leftover = maps.Clone(in)
	_ = c.ExtractValue(leftover)
	out = url.Values{}
	c.AppendValue(out)
	return out, leftover
}

func assertValuesEq(t *testing.T, context string, want, have url.Values) {
	t.Helper()
	assert.MapsEqFunc(t, context, want, have, slices.Equal[[]string])
}

func abSelect(multiple bool) *hmc.Select {
	return &hmc.Select{Name: "s", Multiple: multiple, Options: []hmc.Option{
		{Selected: true}, {Value: "a"}, {Value: "b"},
	}}
}

// Lossless round trips: the appended values equal the submitted form,
// and the control claims everything.

func TestRoundTripInputText(t *testing.T) {
	form := url.Values{"i": {"x"}}
	out, leftover := roundTrip(&hmc.Input{Name: "i"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripInputAbsent(t *testing.T) {
	form := url.Values{}
	out, leftover := roundTrip(&hmc.Input{Name: "i"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripCheckboxChecked(t *testing.T) {
	form := url.Values{"c": {"on"}}
	out, leftover := roundTrip(&hmc.Input{Name: "c", Type: "checkbox", Value: "on"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripCheckboxUnchecked(t *testing.T) {
	form := url.Values{}
	out, leftover := roundTrip(&hmc.Input{Name: "c", Type: "checkbox", Value: "on"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripCheckboxEmptyValue(t *testing.T) {
	form := url.Values{"c": {""}}
	out, leftover := roundTrip(&hmc.Input{Name: "c", Type: "checkbox"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripSelectListed(t *testing.T) {
	form := url.Values{"s": {"b"}}
	out, leftover := roundTrip(abSelect(false), form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripSelectDefault(t *testing.T) {
	form := url.Values{}
	out, leftover := roundTrip(abSelect(false), form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripSelectUnlisted(t *testing.T) {
	form := url.Values{"s": {"z"}}
	out, leftover := roundTrip(abSelect(false), form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripMultiselectInOptionOrder(t *testing.T) {
	form := url.Values{"s": {"a", "b"}}
	out, leftover := roundTrip(abSelect(true), form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripNamedMap(t *testing.T) {
	form := url.Values{"d[k]": {"1"}, "d[j]": {"2", "3"}}
	out, leftover := roundTrip(&hmc.Map{Name: "d"}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

func TestRoundTripCatchAllMap(t *testing.T) {
	form := url.Values{"x": {"1"}, "y": {"2"}}
	out, leftover := roundTrip(&hmc.Map{}, form)
	assertValuesEq(t, "appended", form, out)
	assert.Eq(t, "leftover entries", 0, len(leftover))
}

// Lossy round trips: each test here is a deliberate loss; changing one
// should be a decision.

// An empty value means "unset", so it isn't echoed back.
func TestRoundTripLossyInputEmpty(t *testing.T) {
	out, leftover := roundTrip(&hmc.Input{Name: "i"}, url.Values{"i": {""}})
	assertValuesEq(t, "appended", url.Values{}, out)
	assertValuesEq(t, "leftover", url.Values{}, leftover)
}

// An empty value means "unset", so it isn't echoed back.
func TestRoundTripLossySelectEmpty(t *testing.T) {
	out, leftover := roundTrip(abSelect(false), url.Values{"s": {""}})
	assertValuesEq(t, "appended", url.Values{}, out)
	assertValuesEq(t, "leftover", url.Values{}, leftover)
}

// Single-valued controls claim only the first value.
func TestRoundTripLossyInputRepeated(t *testing.T) {
	out, leftover := roundTrip(&hmc.Input{Name: "i"}, url.Values{"i": {"x", "y"}})
	assertValuesEq(t, "appended", url.Values{"i": {"x"}}, out)
	assertValuesEq(t, "leftover", url.Values{"i": {"y"}}, leftover)
}

// Single-valued controls claim only the first value.
func TestRoundTripLossySelectRepeated(t *testing.T) {
	out, leftover := roundTrip(abSelect(false), url.Values{"s": {"a", "b"}})
	assertValuesEq(t, "appended", url.Values{"s": {"a"}}, out)
	assertValuesEq(t, "leftover", url.Values{"s": {"b"}}, leftover)
}

// Order follows Options, not the submission.
func TestRoundTripLossyMultiselectReordered(t *testing.T) {
	out, leftover := roundTrip(abSelect(true), url.Values{"s": {"b", "a"}})
	assertValuesEq(t, "appended", url.Values{"s": {"a", "b"}}, out)
	assertValuesEq(t, "leftover", url.Values{}, leftover)
}

// Unlisted values are prepended one by one, so they come back reversed.
func TestRoundTripLossyMultiselectUnlisted(t *testing.T) {
	out, leftover := roundTrip(abSelect(true), url.Values{"s": {"y", "z"}})
	assertValuesEq(t, "appended", url.Values{"s": {"z", "y"}}, out)
	assertValuesEq(t, "leftover", url.Values{}, leftover)
}

// Disabled controls neither read nor write.
func TestRoundTripLossyInputDisabled(t *testing.T) {
	out, leftover := roundTrip(&hmc.Input{Name: "i", Disabled: true}, url.Values{"i": {"x"}})
	assertValuesEq(t, "appended", url.Values{}, out)
	assertValuesEq(t, "leftover", url.Values{"i": {"x"}}, leftover)
}

// Disabled controls neither read nor write.
func TestRoundTripLossySelectDisabled(t *testing.T) {
	out, leftover := roundTrip(&hmc.Select{Name: "s", Disabled: true}, url.Values{"s": {"a"}})
	assertValuesEq(t, "appended", url.Values{}, out)
	assertValuesEq(t, "leftover", url.Values{"s": {"a"}}, leftover)
}
