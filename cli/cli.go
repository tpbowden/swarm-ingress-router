package cli

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/tpbowden/swarm-ingress-router/types"
)

func GetConfig(args []string) (types.Startable, types.Configuration) {
	var c types.Configuration
	var s types.Startable
	envconfig.Process("ingress", &c)
	return s, c
}
