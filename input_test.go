package hmc_test

import (
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

func TestInputMin(t *testing.T) {
	i := hmc.Input{
		Min: "123",
	}
	i.Validate()
	assert.Eq(t, "no error", "", i.Error)
}
