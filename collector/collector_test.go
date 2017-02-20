package collector

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/store"
	"github.com/tpbowden/swarm-ingress-router/types"
)

var config = types.Configuration{}

var storeContent = map[string]string{}

var swarmServices = []swarm.Service{
	{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "some service",
			},
		},
	},
}

var parsedServices = []types.Service{
	{
		Name: "some service",
	},
}

var serializedServices = "some parsed services"

var servicesForParsing []swarm.Service
var servicesForSerialization []types.Service

type fakeStore struct {
}

func (f fakeStore) Set(key, value string) {
	storeContent[key] = value
}

func fakeServiceList() []swarm.Service {
	return swarmServices
}

func fakeParser(services []swarm.Service) []types.Service {
	servicesForParsing = services
	return parsedServices
}

func fakeSerializer(services []types.Service) string {
	servicesForSerialization = services
	return serializedServices
}

var subject = collector{
	config:            config,
	serviceList:       fakeServiceList,
	parseServices:     fakeParser,
	serializeServices: fakeSerializer,
	store:             store.Store(fakeStore{}),
}

func TestRunningTheCollector(t *testing.T) {
	subject.update()

	if servicesForParsing[0].Spec.Annotations.Name != swarmServices[0].Spec.Annotations.Name {
		t.Error("Swarm services were not sent for parsing")
	}

	if servicesForSerialization[0] != parsedServices[0] {
		t.Error("Parsed services were not sent for serialization")
	}

	if storeContent["services"] != serializedServices {
		t.Error("Store did not set services correctly")
	}
}
