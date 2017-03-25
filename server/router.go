package server

import (
	"github.com/tpbowden/swarm-ingress-router/service"
)

type hostname string

type router struct {
	services map[hostname]service.Service
}

func (r *router) subscribe() {

}

func route(request) service.Service {
	var s service.Service
	return s
}

func newRouter() *router {
	var services map[hostname]service.Service

	return &router{
		services: services,
	}
}
