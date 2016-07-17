package service

import (
	"crypto/tls"
	"fmt"
)

type TLSService struct {
	name        string
	port        int
	dnsName     string
	certificate tls.Certificate
	forceTLS    bool
}

func (s TLSService) Certificate() (*tls.Certificate, bool) {
	return &s.certificate, true
}

func (s TLSService) DNSName() string {
	return s.dnsName
}

func (s TLSService) URL() string {
	return fmt.Sprintf("http://%s:%d", s.name, s.port)
}

func (s TLSService) ForceTLS() bool {
	return s.forceTLS
}

func NewTLSService(name string, port int, dnsName string, certificate tls.Certificate, forceTLS bool) TLSService {
	return TLSService{name: name, port: port, dnsName: dnsName, certificate: certificate, forceTLS: forceTLS}
}
