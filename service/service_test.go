package service

import (
  "testing"
)

func TestServiceURL(t *testing.T) {
  s := NewService("foo", 8080, "bar")
  url := s.Url()
  expected := "http://foo:8080"

  if url != expected {
    t.Error("Expected", expected, "got", url)
  }
}
