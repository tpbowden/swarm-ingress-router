package service

import (
	"testing"
)

func TestServiceURL(t *testing.T) {
	s := NewService("foo", 8080, "bar")
	url := s.URL()
	expected := "http://foo:8080"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestServiceDNSNAme(t *testing.T) {
	s := Service{dnsName: "something"}

	expected := "something"
	actual := s.DNSName()

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestServiceForceTLS(t *testing.T) {
	s := Service{forceTLS: true}

	expected := true
	actual := s.ForceTLS()

	if expected != actual {
		t.Errorf("Expected %b, got %b", expected, actual)
	}
}
