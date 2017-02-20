package collector

import (
	"log"
	"time"

	"github.com/docker/docker/api/types/swarm"

	"github.com/tpbowden/swarm-ingress-router/store"
	"github.com/tpbowden/swarm-ingress-router/types"
)

type collector struct {
	config            types.Configuration
	serviceList       func() []swarm.Service
	parseServices     func([]swarm.Service) []types.Service
	serializeServices func([]types.Service) string
	store             store.Store
}

func (c collector) update() {
	log.Print("Starting collector run")

	defer func() {
		if r := recover(); r != nil {
			log.Print("Collector error - ", r)
		}
	}()

	swarmServices := c.serviceList()
	parsedServices := c.parseServices(swarmServices)
	json := c.serializeServices(parsedServices)
	c.store.Set("services", json)

	log.Print("Collector run completed")
}

func (c collector) Start() {
	c.update()

	for range time.Tick(c.config.PollInterval) {
		c.update()
	}
}

func NewCollector(config types.Configuration) types.Startable {
	return collector{
		config:            config,
		serviceList:       newDockerClient().serviceList,
		parseServices:     parseServices,
		serializeServices: serializeServices,
		store:             store.NewStore(config),
	}
}
