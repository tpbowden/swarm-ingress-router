package collector

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/service"
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

var parsedServices = []service.Service{
	{
		Name: "some service",
	},
}

var serializedServices = "some parsed services"

var servicesForParsing []swarm.Service
var servicesForSerialization []service.Service

type fakeStore struct {
}

func (f fakeStore) Set(key, value string) {
	storeContent[key] = value
}

func fakeServiceList() []swarm.Service {
	return swarmServices
}

func fakeParser(services []swarm.Service) []service.Service {
	servicesForParsing = services
	return parsedServices
}

func fakeSerializer(services []service.Service) string {
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

	if !reflect.DeepEqual(servicesForParsing, swarmServices) {
		t.Error("Swarm services were not sent for parsing")
	}

	if !reflect.DeepEqual(servicesForSerialization, parsedServices) {
		t.Error("Parsed services were not sent for serialization")
	}

	if storeContent["services"] != serializedServices {
		t.Error("Store did not set services correctly")
	}
}
