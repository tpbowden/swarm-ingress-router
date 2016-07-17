package service

import (
	"testing"
)

func TestTLSServiceForceTLS(t *testing.T) {
	s := TLSService{forceTLS: true}

	expected := true
	actual := s.ForceTLS()

	if expected != actual {
		t.Errorf("Expected %b, got %b", expected, actual)
	}
}
