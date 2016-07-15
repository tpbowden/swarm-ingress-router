package service

import (
	"log"
	"strconv"

  "github.com/docker/engine-api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/docker"
	"github.com/tpbowden/swarm-ingress-router/router"
)

func LoadAll(client docker.ServicePuller) []router.Routable {
  filters := map[string]string{"label": "ingress=true"}

	services := client.GetServices(filters)
  return parseServices(services)
}

func parseServices(services []swarm.Service) []router.Routable {
	var serviceList []router.Routable

	for _, s := range services {
		port, err := strconv.Atoi(s.Spec.Annotations.Labels["ingress.targetport"])
		if err != nil {
			log.Printf("Invalid port detected for service %s", s.Spec.Annotations.Name)
      continue
    }

    parsedService := NewService(s.Spec.Annotations.Name, port, s.Spec.Annotations.Labels["ingress.dnsname"])
    serviceList = append(serviceList, router.Routable(parsedService))
	}

	return serviceList
}
