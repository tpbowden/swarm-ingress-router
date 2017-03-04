package service

import (
	"reflect"
	"testing"
)

type serializeTest struct {
	description string
	services    []Service
	result      string
}

var tests = []serializeTest{
	{
		description: "Serializing a service",
		services: []Service{
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
		result := Serialize(test.services)

		if result != test.result {
			t.Errorf("Result does not match: expected %v, got %v", test.result, result)
		}

		deserialized := Deserialize(result)

		if !reflect.DeepEqual(deserialized, test.services) {
			t.Errorf("Failed to deserialize services: expected %v, got %v", test.services, deserialized)
		}
	}
}
