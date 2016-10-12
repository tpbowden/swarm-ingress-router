package docker

import (
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// HTTPClient holds the configration needed for a connection to Docker's API
type HTTPClient struct {
	cli *client.Client
}

// Client takes labels and returns matching docker services
type Client interface {
	GetServices(map[string]string) []swarm.Service
}

// GetServices returns all Docker services mathcing the labels giben
func (c *HTTPClient) GetServices(filterList map[string]string) []swarm.Service {
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to lookup services: ", r)
		}
	}()

	filter := filters.NewArgs()
	for k, v := range filterList {
		filter.Add(k, v)
	}

	services, err := c.cli.ServiceList(context.Background(), types.ServiceListOptions{Filter: filter})
	if err != nil {
		log.Print("Failed to lookup services: ", err)
		return []swarm.Service{}
	}

	return services
}

// NewClient returns a new instance of the HTTP client
func NewClient() Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)

	if err != nil {
		log.Panic("Failed to create Docker client: ", err)
	}

	return Client(&HTTPClient{
		cli: cli,
	})
}
