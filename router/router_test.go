package router

import (
  "testing"
)

type TestService struct {}

func (t TestService) DnsName() string {
  return "www.example.com"
}

func (t TestService) Url() string {
  return "http://route.local"
}

func TestAddingARoute(t *testing.T) {
  r := NewRouter()
  service := Routable(TestService{})
  services := make([]Routable, 1)
  services[0] = service
  r.UpdateTable(services)
  routedService, ok := r.Route("www.example.com")

  if !ok {
    t.Fatal("Failed to lookup route")
  }

  actual := routedService.Url()
  expected := "http://route.local"

  if expected != actual {
    t.Error("Expected service URL to equal %s, got %s", expected, actual)
  }
}
