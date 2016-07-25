package service

import (
	"crypto/tls"
	"fmt"
)

// Service holds all metdata required to router to a Docker service
type Service struct {
	Name        string
	Port        int
	DNSName     string
	Secure      bool
	ForceTLS    bool
	EncodedCert string
	EncodedKey  string
}

// Certificate returns a parsed TLS certificate/key pair
func (s Service) Certificate() (tls.Certificate, error) {
	return tls.X509KeyPair([]byte(s.EncodedCert), []byte(s.EncodedKey))
}

// URL returns the URL for a service as a string
func (s Service) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Name, s.Port)
}

// NewService returns a new service instance
func NewService(name string, port int, dnsName string, secure bool, forceTLS bool, encodedCert string, encodedKey string) Service {
	return Service{
		Name:        name,
		Port:        port,
		DNSName:     dnsName,
		Secure:      secure,
		ForceTLS:    forceTLS,
		EncodedCert: encodedCert,
		EncodedKey:  encodedKey,
	}
}
