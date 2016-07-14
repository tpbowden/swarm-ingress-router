package docker

import (
	"log"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"
)

type ServicePuller interface {
	GetServices() []swarm.Service
}

type Client struct {
	socket     string
	apiVersion string
}

func (c Client) GetServices() []swarm.Service {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to lookup services: ", r)
		}
	}()

	if err != nil {
		log.Print("Failed to lookup services: ", err)
		return make([]swarm.Service, 0)
	}

	filter := filters.NewArgs()

	filter.Add("label", "ingress=true")

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filter: filter})
	if err != nil {
		panic(err)
	}

	return services
}

func NewClient() ServicePuller {
	return ServicePuller(Client{socket: "unix:///var/run/docker.sock", apiVersion: "v1.24"})
}
