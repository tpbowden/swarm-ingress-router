package router

import (
  "log"
  "github.com/tpbowden/ingress-router/service"
)

type Router struct {
  routes map[string]*service.Service
}

func (r *Router) Route(address string) (*service.Service, bool) {
  route, ok := r.routes[address]
  return route, ok
}

func (r *Router) UpdateTable(services []service.Service) {
  newTable := make(map[string]*service.Service)
  for _, s := range services {
    log.Printf("Registering service for %s", s.DnsName)
    newTable[s.DnsName] = &s
  }

  r.routes = newTable
}

func NewRouter() *Router {
  return new(Router)
}
