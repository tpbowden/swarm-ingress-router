package router

import (
	"crypto/tls"
	"testing"
)

type TestService struct {
	dnsName string
	url     string
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

	actual, ok := r.RouteToService("www.example.com")
	expected := "http://route.local"

	if !ok {
		t.Error("Failed to lookup %www.example.com")
	}

	if actual.String() != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestRoutingToInvalidMissing(t *testing.T) {
	r := NewRouter()
	service := Routable(TestService{dnsName: "", url: ""})
	services := make([]Routable, 1)
	services[0] = service
	r.UpdateTable(services)

	_, ok := r.RouteToService("www.nowhere.com")

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

	_, ok := r.RouteToService("www.nowhere.com")

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
