package service

import (
  "testing"
)

func TestServiceURL(t *testing.T) {
  s := Service{Name: "foo", Port: 8080, DnsName: "bar"}
  url := s.Url()
  expected := "http://foo:8080"

  if url != expected {
    t.Error("Expected", expected, "got", url)
  }
}
