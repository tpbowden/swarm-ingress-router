package router

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tpbowden/swarm-ingress-router/service"
)

// Router holds the current routing table
type Router struct {
	routes map[string]service.Service
}

// RouteToService returns the correct HTTP handler for a given service's DNS name
func (r *Router) RouteToService(address string, secure bool) (http.Handler, bool) {
	var handler http.Handler

	route, ok := r.routes[address]
	if !ok {
		log.Printf("Failed to lookup service for %s", address)
		return handler, false
	}

	if secure && !route.Secure {
		return handler, false
	}

	serviceURL, err := url.Parse(route.URL())
	if err != nil {
		log.Printf("Failed to parse URL for service %s", address)
		return handler, false
	}

	if secure || !route.ForceTLS {
		return http.Handler(httputil.NewSingleHostReverseProxy(serviceURL)), true
	}

	redirectAddress := fmt.Sprintf("https://%s", address)
	return NewRedirectHandler(redirectAddress, 301), true
}

// CertificateForService returns the certificate for a service (if one exists)
func (r *Router) CertificateForService(address string) (*tls.Certificate, bool) {
	var cert *tls.Certificate

	route, ok := r.routes[address]
	if !ok {
		log.Printf("Failed to lookup service for %s", address)
		return cert, false
	}

	certificate, err := route.Certificate()
	if err != nil {
		log.Print("Certificate parse failure", err)
		return cert, false
	}

	return &certificate, true
}

// UpdateTable is an atomic operation to update the routing table
func (r *Router) UpdateTable(services []service.Service) {
	newTable := make(map[string]service.Service)

	for _, s := range services {
		log.Printf("Registering service for %s", s.DNSName)
		newTable[s.DNSName] = s
	}

	r.routes = newTable
}

// NewRouter returns a new instance of the router
func NewRouter() *Router {
	return &Router{}
}
