package server

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tpbowden/swarm-ingress-router/cache"
	"github.com/tpbowden/swarm-ingress-router/router"
	"github.com/tpbowden/swarm-ingress-router/service"
)

type Startable interface {
	Start()
}

type Server struct {
	bindAddress string
	cache       cache.Cache
	router      *router.Router
}

func (s *Server) syncServices() {
	var services []service.Service
	servicesJson, getErr := s.cache.Get("services")

	if getErr != nil {
		log.Printf("Failed to load servics from cache", getErr)
		return
	}

	err := json.Unmarshal(servicesJson, &services)

	if err != nil {
		log.Print("Failed to sync services", err)
		return
	}

	s.router.UpdateTable(services)
	log.Printf("Routes updated")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dnsName := strings.Split(req.Host, ":")[0]
	secure := req.TLS != nil

	handler, ok := s.router.RouteToService(dnsName, secure)
	if !ok {
		fmt.Fprintf(w, "Failed to look up service")
		return
	}

	handler.ServeHTTP(w, req)
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
	http.ListenAndServe(bind, s)
}

func (s *Server) startHTTPSServer() {
	bind := fmt.Sprintf("%s:8443", s.bindAddress)
	config := &tls.Config{GetCertificate: s.getCertificate}
	listener, _ := tls.Listen("tcp", bind, config)
	tlsServer := http.Server{Handler: s}

	log.Printf("Server listening for HTTPS on https://%s", bind)
	tlsServer.Serve(listener)
}

func (s *Server) Start() {
	go func() {
		s.syncServices()
		s.cache.Subscribe("inress-router", s.syncServices)
	}()
	go s.startHTTPServer()
	go s.startHTTPSServer()
	select {}
}

func NewServer(bind, redis string) Startable {
	router := router.NewRouter()
	cache := cache.NewCache(redis)
	return Startable(&Server{bindAddress: bind, router: router, cache: cache})
}
