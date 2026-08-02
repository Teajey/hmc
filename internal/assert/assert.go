package assert

import (
	"errors"
	"slices"
	"testing"
)

func Eq[C comparable](t *testing.T, context string, expected, actual C) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: %#v != %#v", context, expected, actual)
	}
}

func True(t *testing.T, context string, condition bool) {
	t.Helper()
	if !condition {
		t.Errorf("false: %s", context)
	}
}

func FatalTrue(t *testing.T, context string, condition bool) {
	t.Helper()
	if !condition {
		t.Fatalf("%s", context)
	}
}

func Err(t *testing.T, context string, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: %s", context, err)
	}
}

func FatalErr(t *testing.T, context string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %s", context, err)
	}
}

func FatalErrIs(t *testing.T, context string, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("%s: encountered unexpected error: %s", context, err)
	}
}

func FatalErrAs(t *testing.T, context string, err error, target any) {
	t.Helper()
	if !errors.As(err, target) {
		t.Fatalf("%s: encountered unexpected error: %s", context, err)
	}
}

func SlicesEq[S ~[]E, E comparable](t *testing.T, context string, expected, actual S) {
	t.Helper()
	if !slices.Equal(expected, actual) {
		t.Errorf("%s: %#v != %#v", context, expected, actual)
	}
}
