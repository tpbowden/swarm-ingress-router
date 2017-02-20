package collector

import (
	"testing"

	"github.com/tpbowden/swarm-ingress-router/types"
)

type serializeTest struct {
	description string
	services    []types.Service
	result      string
}

var tests = []serializeTest{
	{
		description: "Serializing a service",
		services: []types.Service{
			{
				Name:        "something",
				DNSNames:    []string{"something.local"},
				Port:        100,
				Secure:      true,
				ForceTLS:    true,
				Certificate: "a cert",
				Key:         "a key",
			},
		},
		result: "[{\"Name\":\"something\",\"DNSNames\":[\"something.local\"],\"Port\":100,\"Certificate\":\"a cert\",\"Key\":\"a key\",\"Secure\":true,\"ForceTLS\":true}]",
	},
}

func TestSerializingServices(t *testing.T) {
	for _, test := range tests {
		result := serializeServices(test.services)

		if result != test.result {
			t.Errorf("Result does not match: expected %v, got %v", test.result, result)
		}
	}
}
