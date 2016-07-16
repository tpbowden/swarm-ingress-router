package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	router       *router.Router
}

func (s *Server) updateServices() {
	log.Print("Updating routes")
	client := docker.NewClient()
	services := service.LoadAll(client)
	s.router.UpdateTable(services)
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
	go s.startTicker()
	go s.startHTTPServer()
	go s.startHTTPSServer()
	select {}

}

func NewServer(bind string, pollInterval int) Startable {
	router := router.NewRouter()
	return Startable(&Server{bindAddress: bind, router: router, pollInterval: time.Duration(pollInterval)})
}
