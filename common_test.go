package gos3headersetter

import (
	"testing"
)

func TestEndsWithTrue(t *testing.T) {
	assertEqualBool(t, "suffix", true, endsWith("ab", "b"))
}

func TestEndsWithFalse(t *testing.T) {
	assertEqualBool(t, "suffix", false, endsWith("ab", "a"))
}

func TestCompareWhenNoUpdate(t *testing.T) {
	change, value := compare("x", "")
	assertEqualBool(t, "change", change, false)
	assertEqualString(t, "value", value, "x")
}

func TestCompareWhenNewValueSame(t *testing.T) {
	change, value := compare("x", "x")
	assertEqualBool(t, "change", change, false)
	assertEqualString(t, "value", value, "x")
}

func TestCompareWhenNewValueDifferent(t *testing.T) {
	change, value := compare("x", "y")
	assertEqualBool(t, "change", change, true)
	assertEqualString(t, "value", value, "y")
}
