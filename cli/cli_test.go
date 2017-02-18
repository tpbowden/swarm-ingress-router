package cli

import (
	"os"
	"testing"
	"time"

	"github.com/tpbowden/swarm-ingress-router/types"
)

func TestDefaults(t *testing.T) {
	defaultValues := types.Configuration{
		Redis:        "localhost:6379",
		Bind:         "0.0.0.0",
		PollInterval: 10 * time.Second,
	}

	_, config := GetConfig([]string{"server"})
	if config != defaultValues {
		t.Errorf("Defaults did not match: expected %+v, got %+v", defaultValues, config)
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

	_, config := GetConfig([]string{"server"})

	if config != overriddenValues {
		t.Errorf("Overrides did not match: expected %+v, got %+v", overriddenValues, config)
	}

}
