package collector

import (
	"log"
	"time"

	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/types"
)

type Collector struct {
	config            types.Configuration
	docker            dockerClient
	parseServices     func([]swarm.Service) []types.Service
	serializeServices func([]types.Service) string
}

func (c Collector) update() {
	defer func() {
		if r := recover(); r != nil {
			log.Print("Collector error - ", r)
		}
	}()

	swarmServices := c.docker.serviceList()
	parsedServices := c.parseServices(swarmServices)
	json := c.serializeServices(parsedServices)
	print(json)
}

func (c Collector) Start() {
	c.update()

	for range time.Tick(c.config.PollInterval) {
		c.update()
	}
}

func NewCollector(config types.Configuration) types.Startable {
	return Collector{
		config:            config,
		docker:            newDockerClient(),
		parseServices:     parseServices,
		serializeServices: serializeServices,
	}
}
