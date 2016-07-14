package router

type Routable interface {
	DNSName() string
	URL() string
}
