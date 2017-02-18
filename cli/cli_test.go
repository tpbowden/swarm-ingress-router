package cli

import (
	"os"
	"testing"
	"time"

	"github.com/tpbowden/swarm-ingress-router/types"
)

var currentConfig types.Configuration
var appType string
var exitMessage string

type fakeApp struct {
}

func (f fakeApp) Start() {
	return
}

func fakeAbort(msg string) {
	exitMessage = msg
}

func fakeServer(config types.Configuration) types.Startable {
	currentConfig = config
	appType = "server"
	return types.Startable(fakeApp{})
}

func fakeCollector(config types.Configuration) types.Startable {
	currentConfig = config
	appType = "collector"
	return types.Startable(fakeApp{})
}

var subject = CLI{
	initServer:    fakeServer,
	initCollector: fakeCollector,
	abort:         fakeAbort,
}

func TestServerDefaults(t *testing.T) {
	defaultValues := types.Configuration{
		Redis:        "localhost:6379",
		Bind:         "0.0.0.0",
		PollInterval: 10 * time.Second,
	}

	subject.GetConfig([]string{"server"})
	if currentConfig != defaultValues {
		t.Errorf("Defaults did not match: expected %+v, got %+v", defaultValues, currentConfig)
	}

	if appType != "server" {
		t.Errorf("Expected app to be a server, got %v", appType)
	}

}

func TestCollectorDefaults(t *testing.T) {
	defaultValues := types.Configuration{
		Redis:        "localhost:6379",
		Bind:         "0.0.0.0",
		PollInterval: 10 * time.Second,
	}

	subject.GetConfig([]string{"collector"})
	if currentConfig != defaultValues {
		t.Errorf("Defaults did not match: expected %+v, got %+v", defaultValues, currentConfig)
	}

	if appType != "collector" {
		t.Errorf("Expected app to be a collector, got %v", appType)
	}
}

func TestOverrides(t *testing.T) {
	overriddenValues := types.Configuration{
		Redis:        "some-address:1234",
		Bind:         "1.2.3.4",
		PollInterval: 5 * time.Minute,
	}

	os.Setenv("INGRESS_REDIS", "some-address:1234")
	os.Setenv("INGRESS_BIND", "1.2.3.4")
	os.Setenv("INGRESS_POLL_INTERVAL", "5m")

	subject.GetConfig([]string{"server"})

	if currentConfig != overriddenValues {
		t.Errorf("Overrides did not match: expected %+v, got %+v", overriddenValues, currentConfig)
	}
}

func TestAbortingOnNoArguments(t *testing.T) {
	subject.GetConfig([]string{})

	if exitMessage != usage {
		t.Error("Expected to have aborted with the usage message")
	}
}

func TestAbortingOnInvalidArguments(t *testing.T) {
	subject.GetConfig([]string{"something"})

	if exitMessage != usage {
		t.Error("Expected to have aborted with the usage message")
	}
}
