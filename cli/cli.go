package cli

import (
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/tpbowden/swarm-ingress-router/collector"
	"github.com/tpbowden/swarm-ingress-router/server"
	"github.com/tpbowden/swarm-ingress-router/types"
)

const usage = "Please specify either 'server' or 'collector' as the first and only argument"

func GetConfig(args []string) (types.Startable, types.Configuration) {
	var c types.Configuration
	var s types.Startable

	if len(args) != 1 {
		log.Fatal(usage)
	}

	command := args[0]
	if command == "server" {
		s = server.NewServer()
	} else if command == "collector" {
		s = collector.NewCollector()
	} else {
		log.Fatal(usage)
	}

	envconfig.Process("ingress", &c)
	return s, c
}
