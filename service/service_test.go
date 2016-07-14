package service

import (
  "testing"
)

func TestServiceURL(t *testing.T) {
  s := NewService("foo", 8080, "bar")
  url := s.Url()
  expected := "http://foo:8080"

  if url != expected {
    t.Errorf("Expected %s, got %s", expected, url)
  }
}

func TestServiceDNSNAme(t *testing.T) {
  s := Service{dnsName: "something"}

  expected := "something"
  actual := s.DnsName()

  if expected != actual {
    t.Errorf("Expected %s, got %s", expected, actual)
  }
}

