package router

import (
	"io/ioutil"
	"net/http/httputil"
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
	RouterTest{
		description: "A valid service returns an HTTP proxy",
		host:        "example.local",
		success:     true,
		proxy:       true,
		services: []service.Service{
			service.Service{
				Name:     "my-service",
				Port:     3000,
				DNSName:  "example.local",
				Secure:   false,
				ForceTLS: false,
			},
		},
	},
	RouterTest{
		description: "A secure connection to an insecure service is not successful",
		host:        "example.local",
		success:     false,
		secure:      true,
		services: []service.Service{
			service.Service{
				Name:    "my-service",
				Port:    3000,
				DNSName: "example.local",
				Secure:  false,
			},
		},
	},
	RouterTest{
		description: "A missing service does not return successfully",
		host:        "example.local",
		services:    []service.Service{},
	},
	RouterTest{
		description: "A service with an invalid URL does not return successfully",
		host:        "example.local",
		services: []service.Service{
			service.Service{
				Name:    "[::1]a",
				Port:    3000,
				DNSName: "example.local",
			},
		},
	},
	RouterTest{
		description: "An insecure connection with forceTLS returns a redirect",
		host:        "example.local",
		success:     true,
		redirect:    true,
		services: []service.Service{
			service.Service{
				Name:     "my-service",
				Port:     3000,
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

		result, ok := subject.RouteToService(test.host, test.secure)

		if ok != test.success {
			t.Errorf("Test failed: service fetching did not match: %s", test.description)
		}

		if test.redirect {
			_, assertOk := result.(*RedirectHandler)
			if !assertOk {
				t.Errorf("Test failed: expected a redirect: %s", test.description)
			}
		}

		if test.proxy {
			_, assertOk := result.(*httputil.ReverseProxy)
			if !assertOk {
				t.Errorf("Test failed: expected a reverse proxy: %s", test.description)
			}
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
	CertificateTest{
		description: "Missing services do not return successfully",
		services: []service.Service{
			service.Service{
				Name:    "my-service",
				Port:    3000,
				DNSName: "example.local",
				Secure:  false,
			},
		},
		host:    "foo.local",
		success: false,
	},
	CertificateTest{
		description: "Invalid certificates do not return successfully",
		services: []service.Service{
			service.Service{
				Name:        "my-service",
				Port:        3000,
				DNSName:     "example.local",
				Secure:      true,
				EncodedCert: "some data",
				EncodedKey:  "some data",
			},
		},
		host:    "example.local",
		success: false,
	},
	CertificateTest{
		description: "Valid certificates return successfully",
		services: []service.Service{
			service.Service{
				Name:        "my-service",
				Port:        3000,
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
