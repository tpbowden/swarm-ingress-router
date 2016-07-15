package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/tpbowden/swarm-ingress-router/docker"
	"github.com/tpbowden/swarm-ingress-router/router"
	"github.com/tpbowden/swarm-ingress-router/service"
)

type Startable interface {
	Start()
}

type Server struct {
	bindAddress  string
	pollInterval time.Duration
	router       router.Router
}

func (s *Server) updateServices() {
	log.Print("Updating routes")
	client := docker.NewClient()
	services := service.LoadAll(client)
	s.router.UpdateTable(services)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Starting request for %s", req.Host)
	dnsName := strings.Split(req.Host, ":")[0]

	srv, ok := s.router.Route(dnsName)
	if !ok {
		fmt.Fprintf(w, "Failed to route to service")
		return
	}

	url, err := url.Parse(srv.URL())
	if err != nil {
		fmt.Fprintf(w, "Failed to parse service URL")
		return
	}

	log.Printf("Routing to %s", srv.URL())
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, req)
}

func (s *Server) startTicker() {
	go s.updateServices()

	ticker := time.NewTicker(s.pollInterval * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.updateServices()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Server) getCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Printf("Looking up certificate for %s", clientHello.ServerName)
	srv, ok := s.router.Route(clientHello.ServerName)
	if !ok {
		log.Printf("Failed to look up service for %s", clientHello.ServerName)
		return &tls.Certificate{}, errors.New("No service for host found")
	}

	tlsService, ok := srv.(service.TLSService)

	if !ok {
		return &tls.Certificate{}, errors.New("No TLS service found")
	}

	return tlsService.Certificate(), nil

}

func (s *Server) Start() {
	s.startTicker()
	bind := fmt.Sprintf("%s:8080", s.bindAddress)
	tlsBind := fmt.Sprintf("%s:8443", s.bindAddress)

	log.Printf("Server listening for http on http://%s", bind)
	go http.ListenAndServe(bind, s)

	config := &tls.Config{GetCertificate: s.getCertificate}
	listener, _ := tls.Listen("tcp", tlsBind, config)
	tlsServer := http.Server{Addr: tlsBind, Handler: s}

	log.Printf("Server listening for https on https://%s", tlsBind)
	tlsServer.Serve(listener)
}

func NewServer(bind string, pollInterval int) Startable {
	router := router.NewRouter()
	return Startable(&Server{bindAddress: bind, router: router, pollInterval: time.Duration(pollInterval)})
}
