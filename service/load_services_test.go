package service

import (
	"testing"

	"github.com/docker/engine-api/types/swarm"
)

type FakeClient struct {
	dnsName string
	port    string
}

func (f FakeClient) GetServices(filters map[string]string) []swarm.Service {
	labels := make(map[string]string)
	labels["ingress.targetport"] = f.port
	labels["ingress.dnsname"] = f.dnsName

	fakeService := swarm.Service{
		ID: "123",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "myservice", Labels: labels},
		},
	}

	ignoredService := swarm.Service{
		ID: "567",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "ignored", Labels: map[string]string{}},
		},
	}

	return []swarm.Service{ignoredService, fakeService}
}

func TestLoadingServices(t *testing.T) {
	result, _ := LoadAll(FakeClient{port: "100", dnsName: "foo.bar.baz"})

	parsedService := result[0]

	expectedName := "foo.bar.baz"
	actualName := parsedService.DNSName()

	if expectedName != actualName {
		t.Errorf("Expected DNS name of %s, got %s", expectedName, actualName)
	}

	expectedURL := "http://myservice:100"
	actualURL := parsedService.URL()

	if expectedURL != actualURL {
		t.Errorf("Expected URL of %s, got %s", expectedURL, actualURL)
	}
}

func TestLoadingInvalidService(t *testing.T) {
	result, _ := LoadAll(FakeClient{dnsName: "foo.bar.baz", port: "abc"})
	if len(result) != 0 {
		t.Errorf("Expected no services to be created")
	}
}
