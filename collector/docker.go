package collector

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type dockerClient struct {
	cli *client.Client
}

func (d dockerClient) serviceList() []swarm.Service {
	filter := filters.NewArgs()
	filter.Add("label", "ingress=true")
	options := types.ServiceListOptions{Filters: filter}

	result, err := d.cli.ServiceList(context.Background(), options)

	if err != nil {
		panic(fmt.Sprintf("Failed to load Docker services: %v", err))
	}

	return result
}

func newDockerClient() dockerClient {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.26", nil, defaultHeaders)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Docker client: %v", err))
	}

	return dockerClient{cli: cli}
}
