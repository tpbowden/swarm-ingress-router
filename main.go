package main

import (
	"os"

	"github.com/tpbowden/swarm-ingress-router/cli"
	"github.com/tpbowden/swarm-ingress-router/server"
)

func main() {
	cli.Start(os.Args, server.NewServer)
}
