package router

import (
  "log"
)

type Router struct {
  routes map[string]Routable
}

func (r *Router) Route(address string) (Routable, bool) {
  route, ok := r.routes[address]
  return route, ok
}

func (r *Router) UpdateTable(services []Routable) {
  newTable := make(map[string]Routable)
  for _, s := range services {
    log.Printf("Registering service for %s", s.DnsName())
    newTable[s.DnsName()] = s
  }

  r.routes = newTable
}

func NewRouter() Router {
  return Router{}
}
