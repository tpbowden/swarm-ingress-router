package router

import (
	"crypto/tls"
)

type Routable interface {
	DNSName() string
	URL() string
	Certificate() (*tls.Certificate, bool)
}
