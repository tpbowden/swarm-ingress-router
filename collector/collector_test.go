package collector

import (
	"testing"

	"github.com/tpbowden/swarm-ingress-router/service"
)

type FakeCache struct {
	key  string
	data string
}

type FakePuller struct {
}

func (p *FakePuller) LoadAll() []service.Service {
	return []service.Service{
		{
			DNSNames: []string{"example.com"},
		},
	}
}

func (f *FakeCache) Set(key, value string) error {
	f.key = key
	f.data = value
	return nil
}

func (f *FakeCache) Get(string) ([]byte, error) {
	return nil, nil
}

func (f *FakeCache) Subscribe(string, func()) error {
	return nil
}

func TestUpdatingServices(t *testing.T) {
	fakeCache := &FakeCache{}
	fakePuller := &FakePuller{}
	subject := Collector{cache: fakeCache, servicePuller: fakePuller}
	subject.updateServices()

	if fakeCache.key != "services" {
		t.Error("Expected cache key to equal services, got %s", fakeCache.key)
	}

	expected := "[{\"URL\":\"\",\"DNSNames\":[\"example.com\"],\"Secure\":false,\"ForceTLS\":false,\"EncodedCert\":\"\",\"EncodedKey\":\"\"}]"
	actual := fakeCache.data

	if expected != actual {
		t.Errorf("Expected cache content to equal %s, got %s", expected, actual)
	}
}
