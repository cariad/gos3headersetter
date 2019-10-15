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
