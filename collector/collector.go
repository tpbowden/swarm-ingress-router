package collector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tpbowden/swarm-ingress-router/cache"
	"github.com/tpbowden/swarm-ingress-router/docker"
	"github.com/tpbowden/swarm-ingress-router/service"
)

type Collector struct {
	pollInterval time.Duration
	cache        cache.Cache
}

func (c *Collector) updateServices() {
	log.Print("Updating routes")
	client := docker.NewClient()
	services := service.LoadAll(client)

	json, err := json.Marshal(services)

	if err != nil {
		log.Print("Failed to encode services as json %v", err)
		return
	}

	if cacheError := c.cache.Set("services", string(json)); cacheError != nil {
		log.Printf("Failed to store services in cache: %v", cacheError)
	}

}

func (c *Collector) Start() {
	c.updateServices()

	for range time.Tick(c.pollInterval * time.Second) {
		c.updateServices()
	}
}

func NewCollector(pollInterval int, redis string) Collector {
	cache := cache.NewCache(redis)
	return Collector{pollInterval: time.Duration(pollInterval), cache: cache}
}
