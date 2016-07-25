package service

import (
	"io/ioutil"
	"testing"
)

var certificate, certErr = ioutil.ReadFile("../fixtures/cert.crt")
var key, keyErr = ioutil.ReadFile("../fixtures/key.key")

func TestServiceURL(t *testing.T) {
	s := NewService("foo", 8080, "bar", false, false, "", "")
	url := s.URL()
	expected := "http://foo:8080"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

type CertificateTest struct {
	description string
	cert        string
	key         string
	success     bool
}

var certificateTests = []CertificateTest{
	{
		description: "Invalid certificate returns false",
		cert:        "a cert",
		key:         "a cert",
		success:     false,
	},
	{
		description: "Valid certificate returns true",
		cert:        string(certificate),
		key:         string(key),
		success:     true,
	},
}

func TestServiceCertificate(t *testing.T) {
	if certErr != nil {
		t.Error("Failed to read certificate fixture: ", certErr)
	}

	if keyErr != nil {
		t.Error("Failed to read key fixture: ", keyErr)
	}

	for _, test := range certificateTests {
		subject := Service{EncodedCert: test.cert, EncodedKey: test.key}

		_, err := subject.Certificate()

		ok := err == nil

		if ok != test.success {
			t.Errorf("Failed: %s", test.description)
		}
	}
}
