package service

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

type FakeClient struct {
	services []swarm.Service
}

func (f FakeClient) GetServices(filters map[string]string) []swarm.Service {
	return f.services
}

type LoadServicesTest struct {
	expected    []Service
	description string
	resultCount int
	services    []swarm.Service
}

var loadServicesTests = []LoadServicesTest{
	{
		description: "It parses and returns valid services",
		resultCount: 1,
		expected: []Service{
			{
				URL:         "myservice:8080",
				DNSName:     "example.com",
				Secure:      true,
				ForceTLS:    true,
				EncodedCert: "a cert",
				EncodedKey:  "a key",
			},
		},
		services: []swarm.Service{
			{
				ID: "123",
				Spec: swarm.ServiceSpec{
					Annotations: swarm.Annotations{
						Name: "myservice",
						Labels: map[string]string{
							"ingress.targetport": "8080",
							"ingress.dnsname":    "example.com",
							"ingress.tls":        "true",
							"ingress.forcetls":   "true",
							"ingress.cert":       "a cert",
							"ingress.key":        "a key",
						},
					},
				},
			},
		},
	},
	{
		description: "Discards services with an invalid port",
		expected:    []Service{},
		services: []swarm.Service{
			{
				ID: "123",
				Spec: swarm.ServiceSpec{
					Annotations: swarm.Annotations{
						Name: "myservice",
						Labels: map[string]string{
							"ingress.targetport": "abc",
							"ingress.dnsname":    "example.com",
						},
					},
				},
			},
		},
	},
}

func TestLoadingServices(t *testing.T) {
	for _, test := range loadServicesTests {
		client := FakeClient{services: test.services}
		subject := DockerPuller{client: client}
		result := subject.LoadAll()

		if !(len(result) == len(test.expected)) {
			t.Errorf("Failed: %s", test.description)
		}

		for i := range result {
			if !reflect.DeepEqual(result[i], test.expected[i]) {
				t.Errorf("Equality failed: %s", test.description)
			}
		}
	}
}
