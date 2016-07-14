package service

import (
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/tpbowden/swarm-ingress-router/router"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func LoadAll() []router.Routable {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to lookup services: ", r)
		}
	}()

	if err != nil {
		log.Print("Failed to lookup services: ", err)
		return make([]router.Routable, 0)
	}

	filter := filters.NewArgs()

	filter.Add("label", "ingress=true")

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filter: filter})
	if err != nil {
		panic(err)
	}

	numServices := len(services)

	serviceList := make([]router.Routable, numServices)

	for i, s := range services {
		port, err := strconv.Atoi(s.Spec.Annotations.Labels["ingress.targetport"])
		if err != nil {
			log.Printf("Invalid port detected for service %s", s.Spec.Annotations.Name)
		} else {
			parsedService := NewService(s.Spec.Annotations.Name, port, s.Spec.Annotations.Labels["ingress.dnsname"])
			serviceList[i] = router.Routable(parsedService)
		}

	}

	return serviceList
}
