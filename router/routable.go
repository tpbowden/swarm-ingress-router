package router

import (
	"crypto/tls"
)

type Routable interface {
	URL() string
	Certificate() (*tls.Certificate, error)
}
