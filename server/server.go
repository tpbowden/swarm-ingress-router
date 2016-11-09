package server

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/tpbowden/swarm-ingress-router/cache"
	"github.com/tpbowden/swarm-ingress-router/router"
	"github.com/tpbowden/swarm-ingress-router/service"
	"github.com/tpbowden/swarm-ingress-router/types"
)

// Server holds all state for routing to services
type Server struct {
	bindAddress string
	maxBodySize int
	cache       cache.Cache
	router      *router.Router
}

func (s *Server) syncServices() {
	var services []service.Service
	servicesJSON, getErr := s.cache.Get("services")

	if getErr != nil {
		log.Printf("Failed to load servics from cache: %v", getErr)
		return
	}

	err := json.Unmarshal(servicesJSON, &services)

	if err != nil {
		log.Print("Failed to sync services", err)
		return
	}

	s.router.UpdateTable(services)
	log.Printf("Routes updated")
}

// ServerHTTP is the default HTTP handler for services
func (s *Server) ServeHTTP(ctx *fasthttp.RequestCtx) {
	dnsName := strings.Split(string(ctx.Host()), ":")[0]
	log.Printf("Started %s \"%s\" for %s using host %s", ctx.Method(), ctx.Path(), ctx.RemoteAddr(), dnsName)

	handler, ok := s.router.RouteToService(dnsName, ctx.IsTLS())
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte("Failed to look up service"))
		return
	}

	handler(ctx)
}

func (s *Server) getCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, ok := s.router.CertificateForService(clientHello.ServerName)
	if !ok {
		return cert, errors.New("Failed to lookup certificate")
	}

	return cert, nil
}

func (s *Server) startHTTPServer() {
	bind := fmt.Sprintf("%s:8080", s.bindAddress)
	log.Printf("Server listening for HTTP on http://%s", bind)
	// http.ListenAndServe(bind, s)

	server := &fasthttp.Server{Handler: s.ServeHTTP, MaxRequestBodySize: s.maxBodySize}
	server.ListenAndServe(bind)
}

func (s *Server) startHTTPSServer() {
	bind := fmt.Sprintf("%s:8443", s.bindAddress)
	config := &tls.Config{GetCertificate: s.getCertificate}
	listener, _ := tls.Listen("tcp", bind, config)
	tlsServer := fasthttp.Server{Handler: s.ServeHTTP, MaxRequestBodySize: s.maxBodySize}

	log.Printf("Server listening for HTTPS on https://%s", bind)
	tlsServer.Serve(listener)
}

// Start start the server and listens for changes to the services
func (s *Server) Start() {
	go func() {
		s.syncServices()
		for {
			err := s.cache.Subscribe("inress-router", s.syncServices)
			log.Printf("Subscription to updates lost, retrying in 10 seconds: %v", err)
			time.Sleep(10 * time.Second)
		}
	}()
	go s.startHTTPServer()
	go s.startHTTPSServer()
	select {}
}

// NewServer returns a new instrance of the server
func NewServer(bind, redis string, maxBodySize int) types.Startable {
	router := router.NewRouter()
	cache := cache.NewCache(redis)
	return types.Startable(&Server{bindAddress: bind, maxBodySize: maxBodySize, router: router, cache: cache})
}
