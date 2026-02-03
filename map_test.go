package hmc_test

import (
	"net/url"
	"testing"

	"github.com/Teajey/hmc"
	"github.com/Teajey/hmc/internal/assert"
)

func TestMapExtractMax(t *testing.T) {
	s := hmc.Map{
		Name:       "foo",
		MaxEntries: 1,
	}

	form := url.Values{
		"foo[bar]": {"two"},
		"foo[baz]": {"three"},
	}

	_ = s.ExtractFormValue(form)

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form empty", 0, len(form))
}

func TestMapExtractMaxLength(t *testing.T) {
	s := hmc.Map{
		Name:      "foo",
		MaxValues: 1,
	}

	form := url.Values{
		"foo[bar]": {"two", "three"},
	}

	_ = s.ExtractFormValue(form)

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form empty", 0, len(form))
}

func TestMapExtractMaxKeysLength(t *testing.T) {
	s := hmc.Map{
		Name:         "foo",
		MaxKeyLength: 256,
	}

	form := url.Values{
		"foo[barbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbazbarbaz]": {"two"},
	}

	_ = s.ExtractFormValue(form)

	assert.SnapshotXml(t, s)
	assert.Eq(t, "form empty", 0, len(form))
}
