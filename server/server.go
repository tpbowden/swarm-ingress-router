package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"github.com/tpbowden/swarm-ingress-router/types"
)

type server struct {
	config types.Configuration
	router *router
}

func startServer(s http.Server) {
	log.Printf("Starting HTTP server on %v", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startTLSServer(s http.Server) {
	log.Printf("Starting TLS server on %v", s.Addr)
	if err := s.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}

func (s server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	host := strings.Split(req.Host, ":")[0]
	log.Printf("Request received for %v", host)
}

func (s server) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Printf("Loading certificate for %v", hello.ServerName)
	var cert tls.Certificate
	return &cert, nil
}

func (s server) Start() {
	httpServer := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &s,
	}

	tlsServer := http.Server{
		Addr:    "127.0.0.1:8443",
		Handler: &s,
		TLSConfig: &tls.Config{
			GetCertificate: s.getCertificate,
		},
	}

	go s.router.subscribe()
	go startServer(httpServer)
	go startTLSServer(tlsServer)

	select {}
}

func NewServer(config types.Configuration) types.Startable {
	return server{
		config: config,
		router: newRouter(),
	}
}
