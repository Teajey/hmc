package hmc_test

import (
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

// so that these snapshots are not affected by other formatting changes
type formatted struct {
	Label     string `json:"-"`
	Type      string
	Name      string `json:"-"`
	Value     string
	Error     string
	Required  bool    `json:",omitempty"`
	Disabled  bool    `json:",omitempty"`
	MinLength uint    `json:",omitempty"`
	MaxLength uint    `json:",omitempty"`
	Step      float32 `json:",omitempty"`
	Min       string  `json:",omitempty"`
	Max       string  `json:",omitempty"`
}

func TestInputValidation(t *testing.T) {
	inputs := []hmc.Input{
		{
			Required: true,
		},
		{
			Min: "123",
		},
		{
			Type:  "number",
			Value: "123",
			Min:   "1000",
		},
		{
			Type:  "number",
			Max:   "123",
			Value: "1000",
		},
		{
			Type:      "number",
			MaxLength: 3,
			Value:     "1234",
		},
		{
			Type:      "number",
			MaxLength: 3,
			Value:     "123",
		},
		{
			Type:      "number",
			MinLength: 5,
			Value:     "1234",
		},
		{
			Type:      "number",
			MinLength: 5,
			Value:     "12345",
		},
		{
			Min:   "abc",
			Value: "def",
		},
		{
			Max:   "abc",
			Value: "def",
		},
		{
			Min:   "def",
			Value: "abc",
		},
		{
			Max:   "def",
			Value: "abc",
		},
		{
			Min:   "abc",
			Value: "abc",
		},
		{
			Max:   "abc",
			Value: "abc",
		},
	}
	formatInputs := make([]formatted, len(inputs))
	for i := range inputs {
		_ = inputs[i].Validate()
		formatInputs[i] = formatted(inputs[i])
	}
	assert.SnapshotJson(t, formatInputs)
}

func TestInputInvalidDate(t *testing.T) {
	input := hmc.Input{
		Value: "abc",
	}
	val, _ := input.ParseValueAsDate()
	assert.True(t, "result is unset", val.IsZero())
	assert.SnapshotJson(t, formatted(input))
}

func TestInputInvalidTime(t *testing.T) {
	input := hmc.Input{
		Value: "abc",
	}
	val, _ := input.ParseValueAsTime()
	assert.True(t, "result is unset", val.IsZero())
	assert.SnapshotJson(t, formatted(input))
}

func TestInputInvalidDatetime(t *testing.T) {
	input := hmc.Input{
		Value: "abc",
	}
	val, _ := input.ParseValueAsDatetime()
	assert.True(t, "result is unset", val.IsZero())
	assert.SnapshotJson(t, formatted(input))
}

func TestInputInvalidDatetimeLocal(t *testing.T) {
	input := hmc.Input{
		Value: "abc",
	}
	val, _ := input.ParseValueAsDatetimeLocal()
	assert.True(t, "result is unset", val.IsZero())
	assert.SnapshotJson(t, formatted(input))
}
