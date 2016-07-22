package service

import (
	"crypto/tls"
	"fmt"
	"log"
)

type Service struct {
	Name        string
	Port        int
	DNSName     string
	Secure      bool
	ForceTLS    bool
	EncodedCert string
	EncodedKey  string
}

func (s Service) Certificate() (tls.Certificate, error) {
	log.Print(s.EncodedCert)
	log.Print(s.EncodedKey)
	return tls.X509KeyPair([]byte(s.EncodedCert), []byte(s.EncodedKey))
}

func (s Service) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Name, s.Port)
}

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
