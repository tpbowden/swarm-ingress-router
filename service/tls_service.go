package service

import (
  "crypto/tls"
  "fmt"
)

type TLSService struct {
  name string
  port int
  dnsName string
  certificate tls.Certificate
  forceTLS bool
}

func (s TLSService) Certificate() *tls.Certificate {
  return &s.certificate
}

func (s TLSService) DNSName() string {
  return s.dnsName
}

func (s TLSService) URL() string {
	return fmt.Sprintf("http://%s:%d", s.name, s.port)

}

func NewTLSService(name string, port int, dnsName string, certificate string, key string) TLSService {

  parsedCert, err := tls.X509KeyPair([]byte(certificate), []byte(key))

  if err != nil {
    panic(err)
  }
  return TLSService{name: name, port: port, dnsName: dnsName, certificate: parsedCert}

}
