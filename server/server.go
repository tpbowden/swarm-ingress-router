package server

import (
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

func (s *Server) handler(w http.ResponseWriter, req *http.Request) {
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

func (s Server) Start() {
	s.startTicker()
	http.HandleFunc("/", s.handler)
	bind := fmt.Sprintf("%s:8080", s.bindAddress)
	log.Printf("Server listening on tcp://%s", bind)
	http.ListenAndServe(fmt.Sprintf("%s:8080", s.bindAddress), nil)
}

func NewServer(bind string, pollInterval int) Startable {
	router := router.NewRouter()
	return Startable(Server{bindAddress: bind, router: router, pollInterval: time.Duration(pollInterval)})
}
