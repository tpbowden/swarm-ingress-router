package version

import (
	"testing"
)

func TestVersionNumbers(t *testing.T) {
	version := NewVersionNumber(1,2,3)
	expected := "1.2.3"
	actual := version.String()

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
