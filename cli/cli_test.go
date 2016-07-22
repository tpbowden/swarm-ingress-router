package cli

import (
	"testing"

	"github.com/tpbowden/swarm-ingress-router/server"
)

var (
	fakeServer    TestServer
	serverStarted = false
)

type TestServer struct {
	bindAddress string
	redis       string
}

func (t TestServer) Start() {
	serverStarted = true
}

func newTestServer(bind string, redis string) server.Startable {
	fakeServer.bindAddress = bind
	fakeServer.redis = redis
	return server.Startable(fakeServer)
}

func TestStartingTheServerWithCLI(t *testing.T) {
	args := []string{"cli", "-r", "redis-url", "server", "-b", "1.2.3.4"}
	Start(args, newTestServer)

	expectedAddr := "1.2.3.4"
	actualAddr := fakeServer.bindAddress
	if expectedAddr != actualAddr {
		t.Errorf("Expected bind address to equal %s, got %s", expectedAddr, actualAddr)
	}

	expectedInterval := "redis-url"
	actualInterval := fakeServer.redis
	if expectedInterval != actualInterval {
		t.Errorf("Expected interval to equal %d, got %d", expectedInterval, actualInterval)
	}

	if !serverStarted {
		t.Errorf("Expected the server to be started, but it was not")
	}
}
