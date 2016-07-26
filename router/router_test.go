package router

import (
	"io/ioutil"
	"testing"

	"github.com/tpbowden/swarm-ingress-router/service"
)

var certificate, certErr = ioutil.ReadFile("../fixtures/cert.crt")
var key, keyErr = ioutil.ReadFile("../fixtures/key.key")

type RouterTest struct {
	description string
	services    []service.Service
	host        string
	secure      bool
	success     bool
	redirect    bool
	proxy       bool
}

var routerTests = []RouterTest{
	{
		description: "A valid service returns an HTTP proxy",
		host:        "example.local",
		success:     true,
		proxy:       true,
		services: []service.Service{
			{
				URL:      "http://my-service:3000",
				DNSName:  "example.local",
				Secure:   false,
				ForceTLS: false,
			},
		},
	},
	{
		description: "A secure connection to an insecure service is not successful",
		host:        "example.local",
		success:     false,
		secure:      true,
		services: []service.Service{
			{
				URL:     "http://my-service:3000",
				DNSName: "example.local",
				Secure:  false,
			},
		},
	},
	{
		description: "A missing service does not return successfully",
		host:        "example.local",
		services:    []service.Service{},
	},
	{
		description: "An insecure connection with forceTLS returns a redirect",
		host:        "example.local",
		success:     true,
		redirect:    true,
		services: []service.Service{
			{
				URL:      "http://my-service:3000",
				DNSName:  "example.local",
				ForceTLS: true,
			},
		},
	},
}

func TestRouting(t *testing.T) {
	for _, test := range routerTests {
		subject := NewRouter()
		subject.UpdateTable(test.services)

		_, ok := subject.RouteToService(test.host, test.secure)

		if ok != test.success {
			t.Errorf("Test failed: service fetching did not match: %s", test.description)
		}
	}
}

type CertificateTest struct {
	description string
	services    []service.Service
	host        string
	success     bool
}

var certificateTests = []CertificateTest{
	{
		description: "Missing services do not return successfully",
		services: []service.Service{
			{
				URL:     "http://my-service:3000",
				DNSName: "example.local",
				Secure:  false,
			},
		},
		host:    "foo.local",
		success: false,
	},
	{
		description: "Valid certificates return successfully",
		services: []service.Service{
			{
				URL:         "http://my-service:3000",
				DNSName:     "example.local",
				Secure:      true,
				EncodedCert: string(certificate),
				EncodedKey:  string(key),
			},
		},
		host:    "example.local",
		success: true,
	},
}

func TestCertificates(t *testing.T) {

	if certErr != nil {
		t.Error("Failed to read certificate fixture: ", certErr)
	}

	if keyErr != nil {
		t.Error("Failed to read key fixture: ", keyErr)
	}

	for _, test := range certificateTests {
		subject := NewRouter()
		subject.UpdateTable(test.services)

		_, ok := subject.CertificateForService(test.host)

		if ok != test.success {
			t.Errorf("Test failed: certificate fetching did not match: %s", test.description)
		}

	}
}
