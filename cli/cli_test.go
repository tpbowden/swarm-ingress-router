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
	interval    int
}

func (t TestServer) Start() {
	serverStarted = true
}

func newTestServer(bind string, interval int) server.Startable {
	fakeServer.bindAddress = bind
	fakeServer.interval = interval
	return server.Startable(fakeServer)
}

func TestStartingTheServerWithCLI(t *testing.T) {
	args := make([]string, 5)

	args[0] = "cli"
	args[1] = "-b"
	args[2] = "1.2.3.4"
	args[3] = "-i"
	args[4] = "100"
	Start(args, newTestServer)

	expectedAddr := "1.2.3.4"
	actualAddr := fakeServer.bindAddress

	if expectedAddr != actualAddr {
		t.Errorf("Expected bind address to equal %s, got %s", expectedAddr, actualAddr)
	}

	expectedInterval := 100
	actualInterval := fakeServer.interval
	if expectedInterval != actualInterval {
		t.Errorf("Expected interval to equal %d, got %d", expectedInterval, actualInterval)
	}

	if !serverStarted {
		t.Errorf("Expected the server to be started, but it was not")
	}
}
