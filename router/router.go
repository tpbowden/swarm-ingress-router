package router

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Router struct {
	routes map[string]Routable
}

func (r *Router) RouteToService(address string, secure bool) (http.Handler, bool) {
	var handler http.Handler
	route, ok := r.routes[address]
	if !ok {
		log.Printf("Failed to lookup service for %s", address)
		return handler, false
	}

	serviceURL, err := url.Parse(route.URL())

	if err != nil {
		log.Printf("Failed to parse URL for service %s", address)
		return handler, false
	}

	if secure || !route.ForceTLS() {
		return http.Handler(httputil.NewSingleHostReverseProxy(serviceURL)), true
	}

	redirectAddress := fmt.Sprintf("https://%s", address)
	return NewRedirectHandler(redirectAddress, 301), true

}

func (r *Router) CertificateForService(address string) (*tls.Certificate, bool) {
	var cert *tls.Certificate

	route, ok := r.routes[address]
	if !ok {
		log.Printf("Failed to lookup service for %s", address)
		return cert, false
	}

	return route.Certificate()
}

func (r *Router) UpdateTable(services []Routable) {
	newTable := make(map[string]Routable)
	for _, s := range services {
		log.Printf("Registering service for %s", s.DNSName())
		newTable[s.DNSName()] = s
	}

	r.routes = newTable
}

func NewRouter() *Router {
	return &Router{}
}
