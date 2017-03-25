package collector

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/service"
)

type parseTest struct {
	description   string
	swarmServices []swarm.Service
	result        []service.Service
}

var parseTests = []parseTest{
	{
		description: "Parsing a valid service",
		swarmServices: []swarm.Service{
			{
				Spec: swarm.ServiceSpec{
					Annotations: swarm.Annotations{
						Name: "test service",
						Labels: map[string]string{
							"ingress.targetport": "100",
							"ingress.tls":        "true",
							"ingress.forcetls":   "true",
							"ingress.cert":       "a certificate",
							"ingress.key":        "a key",
							"ingress.dnsnames":   "foo.bar.com bar.com",
						},
					},
				},
			},
		},
		result: []service.Service{
			{
				Name:        "test service",
				DNSNames:    []string{"foo.bar.com", "bar.com"},
				Port:        100,
				Secure:      true,
				ForceTLS:    true,
				Certificate: "a certificate",
				Key:         "a key",
			},
		},
	},
	{
		description: "Skipping an invalid port number",
		swarmServices: []swarm.Service{
			{
				Spec: swarm.ServiceSpec{
					Annotations: swarm.Annotations{
						Name: "test service",
						Labels: map[string]string{
							"ingress.targetport": "abc",
							"ingress.tls":        "true",
							"ingress.forcetls":   "true",
							"ingress.cert":       "a certificate",
							"ingress.key":        "a key",
						},
					},
				},
			},
		},
		result: []service.Service{},
	},
}

func TestParsingServices(t *testing.T) {
	for _, test := range parseTests {
		parsedServices := parseServices(test.swarmServices)

		for i, res := range parsedServices {
			if !reflect.DeepEqual(test.result[i], res) {
				t.Errorf("Services did not match: expected %v, got %v", test.result[i], res)
			}
		}
	}
}
