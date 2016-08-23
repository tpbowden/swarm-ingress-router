package cli

import (
	"github.com/urfave/cli"

	"github.com/tpbowden/swarm-ingress-router/collector"
	"github.com/tpbowden/swarm-ingress-router/server"
	"github.com/tpbowden/swarm-ingress-router/types"
	"github.com/tpbowden/swarm-ingress-router/version"
)

type CLI struct {
	newServer    func(string, string) types.Startable
	newCollector func(int, string) types.Startable
}

// Start initializes the application as a command line app
func (c *CLI) Start(args []string) {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "collector",
			Usage: "start the collector service",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "interval, i",
					Value: 10,
					Usage: "Poll interval in `seconds`",
				},
			},
			Action: func(ctx *cli.Context) error {
				collector := c.newCollector(ctx.Int("interval"), ctx.GlobalString("redis"))
				collector.Start()
				return nil
			},
		},
		{
			Name:  "server",
			Usage: "start the web server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bind, b",
					Value: "127.0.0.1",
					Usage: "Bind to `address`",
				},
			},
			Action: func(ctx *cli.Context) error {
				server := c.newServer(ctx.String("bind"), ctx.GlobalString("redis"))
				server.Start()
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "redis, r",
			Value: "127.0.0.1:6379",
			Usage: "Redis server `address`",
		},
	}

	app.Name = "Swarm Ingress Router"
	app.Usage = "Route DNS names to Swarm services based on labels"
	app.Version = version.Version.String()

	app.Run(args)
}

func NewCLI() CLI {
	return CLI{newServer: server.NewServer, newCollector: collector.NewCollector}
}
