package hmc_test

import (
	"log"
	"net/url"
	"slices"
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

func TestSelectValues(t *testing.T) {
	s := hmc.Select{
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two", Selected: true},
			{Value: "three", Selected: true},
			{Value: "four", Selected: true, Disabled: true},
		},
	}

	values := slices.Collect(s.Values())
	value := s.Value()

	assert.SlicesEq(t, "expected values", []string{"two", "three"}, values)
	assert.Eq(t, "first value", "two", value)
}

func TestSelectMultipleValues(t *testing.T) {
	s := hmc.Select{
		Multiple: true,
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two", Selected: true},
			{Value: "three", Selected: true},
			{Value: "four", Selected: true, Disabled: true},
		},
	}

	values := slices.Collect(s.Values())
	value := s.Value()

	assert.SlicesEq(t, "expected values", []string{"two", "three"}, values)
	assert.Eq(t, "first value", "two", value)
}

func TestSelectExtract(t *testing.T) {
	s := hmc.Select{
		Name: "myval",
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two"},
			{Value: "three"},
			{Value: "four", Disabled: true},
		},
	}

	form := url.Values{
		"myval": {"two", "three"},
	}

	if err := s.ExtractFormValue(form); err != nil {
		s.Error = err.Error()
	}

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form contains leftovers", 1, len(form))
}

func TestSelectExtractEmpty(t *testing.T) {
	s := hmc.Select{
		Name: "myval",
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two"},
			{Value: "three"},
			{Value: "four", Disabled: true},
		},
	}

	form := url.Values{}

	if err := s.ExtractFormValue(form); err != nil {
		s.Error = err.Error()
	}

	assert.SnapshotXml(t, s)
}

func TestSelectMultipleExtract(t *testing.T) {
	s := hmc.Select{
		Name:     "myval",
		Multiple: true,
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two"},
			{Value: "three"},
			{Value: "four", Disabled: true},
		},
	}

	form := url.Values{
		"myval": {"two", "three"},
	}

	if err := s.ExtractFormValue(form); err != nil {
		s.Error = err.Error()
	}

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form empty", 0, len(form))
}

func TestSelectRequired(t *testing.T) {
	s := hmc.Select{
		Name:     "myval",
		Required: true,
		Options: []hmc.Option{
			{Value: "one"},
			{Value: "two"},
			{Value: "three"},
			{Value: "four", Disabled: true},
		},
	}

	form := url.Values{
		"mynotval": {"two", "three"},
	}

	err := s.ExtractFormValue(form)
	if err != nil {
		s.Error = err.Error()
	}
	err = s.Validate()
	if err != nil {
		log.Println("Yes this line is running")
		s.Error = err.Error()
	}

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form contains leftovers", 1, len(form))
}
