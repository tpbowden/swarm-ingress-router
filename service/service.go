package service

import (
	"fmt"
)

type Service struct {
	name    string
	port    int
	dnsName string
}

func (s Service) DnsName() string {
	return s.dnsName
}

func (s Service) Url() string {
	return fmt.Sprintf("http://%s:%d", s.name, s.port)
}

func NewService(name string, port int, dnsName string) Service {
	return Service{name: name, port: port, dnsName: dnsName}
}
