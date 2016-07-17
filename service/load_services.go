package service

import (
	"crypto/tls"
	"errors"
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
		var parsedService router.Routable

		port, err := strconv.Atoi(s.Spec.Annotations.Labels["ingress.targetport"])
		if err != nil {
			log.Printf("Invalid port detected for service %s", s.Spec.Annotations.Name)
			continue
		}

		if s.Spec.Annotations.Labels["ingress.tls"] == "true" {
			forceTLS := s.Spec.Annotations.Labels["ingress.forcetls"] == "true"
			parsedCert, err := extractCertificate(s.Spec.Annotations.Labels)
			if err != nil {
				continue
			}
			parsedService = router.Routable(NewTLSService(
				s.Spec.Annotations.Name,
				port,
				s.Spec.Annotations.Labels["ingress.dnsname"],
				parsedCert,
				forceTLS,
			))
		} else {
			parsedService = router.Routable(NewService(
				s.Spec.Annotations.Name,
				port,
				s.Spec.Annotations.Labels["ingress.dnsname"],
			))
		}

		serviceList = append(serviceList, parsedService)

	}

	return serviceList
}

func extractCertificate(labels map[string]string) (tls.Certificate, error) {
	encodedCert, certOk := labels["ingress.cert"]
	if !certOk {
		return tls.Certificate{}, errors.New("Could not find a certificate")
	}

	encodedKey, keyOk := labels["ingress.key"]
	if !keyOk {
		return tls.Certificate{}, errors.New("Could not find a key")
	}

	return tls.X509KeyPair([]byte(encodedCert), []byte(encodedKey))
}
