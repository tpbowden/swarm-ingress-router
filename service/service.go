package service

import (
	"crypto/tls"
	"fmt"
)

type Service struct {
	name     string
	port     int
	dnsName  string
	forceTLS bool
}

func (s Service) Certificate() (*tls.Certificate, bool) {
	return &tls.Certificate{}, false
}

func (s Service) DNSName() string {
	return s.dnsName
}

func (s Service) URL() string {
	return fmt.Sprintf("http://%s:%d", s.name, s.port)
}

func (s Service) ForceTLS() bool {
	return s.forceTLS
}

func NewService(name string, port int, dnsName string) Service {
	return Service{name: name, port: port, dnsName: dnsName}
}
