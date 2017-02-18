package cli

import (
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/tpbowden/swarm-ingress-router/collector"
	"github.com/tpbowden/swarm-ingress-router/server"
	"github.com/tpbowden/swarm-ingress-router/types"
)

type CLI struct {
	initServer    func(types.Configuration) types.Startable
	initCollector func(types.Configuration) types.Startable
	abort         func(string)
}

const usage = "Please specify either 'server' or 'collector' as the first and only argument"

func (c CLI) GetConfig(args []string) types.Startable {
	var app types.Startable
	if len(args) != 1 {
		c.abort(usage)
		return app
	}

	var config types.Configuration
	envconfig.Process("ingress", &config)

	command := args[0]
	if command == "server" {
		app = c.initServer(config)
	} else if command == "collector" {
		app = c.initCollector(config)
	} else {
		c.abort(usage)
	}

	return app
}

func fatalAbort(msg string) {
	log.Fatal(msg)
}

func NewCLI() CLI {
	return CLI{
		initServer:    server.NewServer,
		initCollector: collector.NewCollector,
		abort:         fatalAbort,
	}
}
