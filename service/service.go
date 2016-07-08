package service

import (
  "fmt"
  "log"
  "strconv"
  "github.com/docker/engine-api/client"
  "github.com/docker/engine-api/types"
  "github.com/docker/engine-api/types/filters"
  "golang.org/x/net/context"
)

type Service struct {
  Name string
  Port int
  DnsName string
}

func (s *Service) Url() string {
  return fmt.Sprintf("http://%s:%d", s.Name, s.Port)
}

func NewService(name string, port int, dnsName string) Service {
  return Service{Name: name, Port: port, DnsName: dnsName}
}

func LoadAll() []Service {
  defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
  cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
  if err != nil {
    panic(err)
  }

  filter := filters.NewArgs()

  filter.Add("label", "ingress=true")

  services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filter: filter})
  if err != nil {
    panic(err)
  }

  numServices := len(services)

  serviceList := make([]Service, numServices)

  for i, s := range services {
    port, err := strconv.Atoi(s.Spec.Annotations.Labels["ingress.targetport"])
    if err != nil {
      log.Printf("Invalid port detected for service %s", s.Spec.Annotations.Name)
    } else {
      parsedService := NewService(s.Spec.Annotations.Name, port, s.Spec.Annotations.Labels["ingress.dnsname"])
      serviceList[i] = parsedService
    }

  }

  return serviceList
}
