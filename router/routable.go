package router

type Routable interface {
  DnsName() string
  Url() string
}

