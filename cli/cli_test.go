package cli

import (
	"testing"

	"github.com/tpbowden/swarm-ingress-router/types"
)

var (
	fakeServer    = &TestServer{}
	fakeCollector = &TestCollector{}
)

type TestServer struct {
	bindAddress string
	redis       string
	maxBodySize int
  readBufferSize int
	started     bool
  cert        string
  key         string
}

func (s *TestServer) Start() {
	s.started = true
}

type TestCollector struct {
	interval int
	redis    string
	started  bool
}

func (c *TestCollector) Start() {
	c.started = true
}

func newTestServer(bind, redis, cert, key string, maxBodySize int, readBufferSize int) types.Startable {
	fakeServer.bindAddress = bind
	fakeServer.redis = redis
  fakeServer.cert = cert
  fakeServer.key = key
	fakeServer.maxBodySize = maxBodySize
  fakeServer.readBufferSize = readBufferSize
	fakeServer.started = false
	return types.Startable(fakeServer)
}

func newTestCollector(interval int, redis string) types.Startable {
	fakeCollector.interval = interval
	fakeCollector.redis = redis
	fakeCollector.started = false
	return types.Startable(fakeCollector)
}

func TestStartingTheServerWithCLI(t *testing.T) {
	args := []string{"cli", "-r", "redis-url", "server", "-c", "wildcard_cert", "-k", "wildcard_cert_key", "-b", "1.2.3.4", "--max-body-size", "412", "--read-buffer-size", "206"}
	subject := CLI{newServer: newTestServer, newCollector: newTestCollector}
	subject.Start(args)

	expectedAddr := "1.2.3.4"
	actualAddr := fakeServer.bindAddress
	if expectedAddr != actualAddr {
		t.Errorf("Expected bind address to equal %s, got %s", expectedAddr, actualAddr)
	}

	expectedRedis := "redis-url"
	actualRedis := fakeServer.redis
	if expectedRedis != actualRedis {
		t.Errorf("Expected redis URL to equal %s, got %s", expectedRedis, actualRedis)
	}

  expectedCert := "wildcard_cert"
  actualCert := fakeServer.cert
  if expectedCert != actualCert {
    t.Errorf("Expected cert to equal %s, got %s", expectedCert, actualCert)
  }

  expectedKey := "wildcard_cert_key"
  actualKey := fakeServer.key
  if expectedKey != actualKey {
    t.Errorf("Expected key to equal %s, got %s", expectedKey, actualKey)
  }

	expectedMaxBodySize := 412 * 1024 * 1024
	actualMaxBodySize := fakeServer.maxBodySize
	if expectedMaxBodySize != actualMaxBodySize {
		t.Errorf("Expected body size limit to equal %d, got %d", expectedMaxBodySize, actualMaxBodySize)
	}

  expectedReadBufferSize := 206 * 1024
	actualReadBufferSize := fakeServer.readBufferSize
	if expectedReadBufferSize != actualReadBufferSize {
		t.Errorf("Expected read buffer size limit to equal %d, got %d", expectedReadBufferSize, actualReadBufferSize)
	}

	if !fakeServer.started {
		t.Errorf("Expected the server to be started, but it was not")
	}
}

func TestStartingTheCollectorWithCLI(t *testing.T) {
	args := []string{"cli", "-r", "redis-url", "collector", "-i", "100"}
	subject := CLI{newServer: newTestServer, newCollector: newTestCollector}
	subject.Start(args)

	expectedInterval := 100
	actualInterval := fakeCollector.interval
	if expectedInterval != actualInterval {
		t.Errorf("Expected interval to equal %i, got %i", expectedInterval, actualInterval)
	}

	expectedRedis := "redis-url"
	actualRedis := fakeCollector.redis
	if expectedRedis != actualRedis {
		t.Errorf("Expected redis URL to equal %s, got %s", expectedRedis, actualRedis)
	}

	if !fakeCollector.started {
		t.Errorf("Expected the collector to be started, but it was not")
	}
}
