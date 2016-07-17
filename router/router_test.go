package router

import (
	"crypto/tls"
	"net/http/httputil"
	"testing"
)

type TestService struct {
	dnsName  string
	url      string
	forceTLS bool
}

func (t TestService) ForceTLS() bool {
	return t.forceTLS
}

func (t TestService) DNSName() string {
	return t.dnsName
}

func (t TestService) URL() string {
	return t.url
}

func (t TestService) Certificate() (*tls.Certificate, bool) {
	return &tls.Certificate{}, true
}

func TestAddingARoute(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "www.example.com", url: "http://route.local"})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	actual, ok := r.RouteToService("www.example.com", true)

	if !ok {
		t.Error("Failed to lookup %www.example.com")
	}

	_, castOk := actual.(*httputil.ReverseProxy)

	if !castOk {
		t.Error("Expected result to be a reverse proxy to the service")
	}
}

func TestRedirectingToHTTPS(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "www.example.com", forceTLS: true, url: "http://route.local"})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	actual, ok := r.RouteToService("www.example.com", false)

	if !ok {
		t.Error("Failed to lookup %www.example.com")
	}

	_, castOk := actual.(*RedirectHandler)

	if !castOk {
		t.Error("Expected result to be a reverse proxy to the service")
	}
}

func TestRoutingToInvalidMissing(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "", url: ""})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	_, ok := r.RouteToService("www.nowhere.com", true)

	if ok {
		t.Error("Expected to fail to route to non-existant service")
	}
}

func TestRoutingToInvalidService(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "www.nowhere.com", url: "http://[::1]a"})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	_, ok := r.RouteToService("www.nowhere.com", true)

	if ok {
		t.Error("Expected to fail to route to invalid service")
	}
}

func TestGettingACertificate(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "www.nowhere.com", url: "http://[::1]a"})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	_, ok := r.CertificateForService("www.nowhere.com")

	if !ok {
		t.Error("Failed to retrieve certificate for service")
	}
}

func TestGettingACertificateForMissingService(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "www.nowhere.com", url: "http://[::1]a"})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	_, ok := r.CertificateForService("www.wrong.com")

	if ok {
		t.Error("Expected to fail to retrieve missing certificate")
	}
}
