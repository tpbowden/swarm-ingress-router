package service

import (
	"log"
	"strconv"

	"github.com/docker/engine-api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/docker"
)

func LoadAll(client docker.ServicePuller) []Service {
	filters := map[string]string{"label": "ingress=true"}

	services := client.GetServices(filters)
	return parseServices(services)
}

func parseServices(services []swarm.Service) []Service {
	var serviceList []Service

	for _, s := range services {
		var parsedService Service

		port, err := strconv.Atoi(s.Spec.Annotations.Labels["ingress.targetport"])
		if err != nil {
			log.Printf("Invalid port detected for service %s", s.Spec.Annotations.Name)
			continue
		}

		secure := s.Spec.Annotations.Labels["ingress.tls"] == "true"
		forceTLS := s.Spec.Annotations.Labels["ingress.forcetls"] == "true"

		parsedService = NewService(
			s.Spec.Annotations.Name,
			port,
			s.Spec.Annotations.Labels["ingress.dnsname"],
			secure,
			forceTLS,
			s.Spec.Annotations.Labels["ingress.cert"],
			s.Spec.Annotations.Labels["ingress.cert"],
		)

		serviceList = append(serviceList, parsedService)

	}

	return serviceList
}
