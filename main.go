package main

import (
	"os"

	"github.com/tpbowden/swarm-ingress-router/cli"
)

func main() {
	app := cli.NewCLI()
	app.Start(os.Args)
}
