package collector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tpbowden/swarm-ingress-router/cache"
	"github.com/tpbowden/swarm-ingress-router/docker"
	"github.com/tpbowden/swarm-ingress-router/service"
	"github.com/tpbowden/swarm-ingress-router/types"
)

// Collector holds all state for the sollector
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
		log.Printf("Failed to encode services as json %v", err)
		return
	}

	if cacheError := c.cache.Set("services", string(json)); cacheError != nil {
		log.Printf("Failed to store services in cache: %v", cacheError)
	}

}

// Start causes the collector to begin polling docker
func (c *Collector) Start() {
	c.updateServices()

	for range time.Tick(c.pollInterval * time.Second) {
		c.updateServices()
	}
}

// NewCollector returns a new instance of the collector
func NewCollector(pollInterval int, redis string) types.Startable {
	cache := cache.NewCache(redis)
	return types.Startable(&Collector{pollInterval: time.Duration(pollInterval), cache: cache})
}
