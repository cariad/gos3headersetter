package gos3headersetter

import (
	"testing"
)

func assertEqualBool(t *testing.T, name string, actual bool, expect bool) {
	if actual != expect {
		t.Errorf("expected %s to be %v but is %v", name, expect, actual)
	}
}

func assertEqualString(t *testing.T, name string, actual string, expect string) {
	if actual != expect {
		t.Errorf("expected %s to be \"%s\" but is \"%s\"", name, expect, actual)
	}
}
