package service

import (
  "crypto/tls"
)

type TLSService struct {
  Service
  certificate tls.Certificate
  forceTLS bool
}

func NewTLSService(name string, port int, dnsName string, certificate string, key string) TLSService {

  return TLSService{}

}
