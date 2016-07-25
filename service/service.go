package service

import (
	"crypto/tls"
	"fmt"
	"log"
)

// Service holds all metdata required to router to a Docker service
type Service struct {
	URL         string
	DNSName     string
	Secure      bool
	ForceTLS    bool
	EncodedCert string
	EncodedKey  string
	parsedCert  tls.Certificate
}

// Certificate returns a parsed TLS certificate/key pair
func (s Service) Certificate() tls.Certificate {
	return s.parsedCert
}

// ParseCertificate parsed the encoded cert / key and stores them on the service
func (s *Service) ParseCertificate() bool {
	if !s.Secure {
		return false
	}

	parsedCert, err := tls.X509KeyPair([]byte(s.EncodedCert), []byte(s.EncodedKey))
	if err != nil {
		log.Printf("Failed to parse certificate for %s", s.DNSName)
		return false
	}

	s.parsedCert = parsedCert
	return true
}

// NewService returns a new service instance
func NewService(name string, port int, dnsName string, secure bool, forceTLS bool, encodedCert string, encodedKey string) Service {
	url := fmt.Sprintf("http://%s:%d", name, port)
	return Service{
		URL:         url,
		DNSName:     dnsName,
		Secure:      secure,
		ForceTLS:    forceTLS,
		EncodedCert: encodedCert,
		EncodedKey:  encodedKey,
	}
}
