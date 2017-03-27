package cli

import (
	"github.com/urfave/cli"

	"github.com/tpbowden/swarm-ingress-router/collector"
	"github.com/tpbowden/swarm-ingress-router/server"
	"github.com/tpbowden/swarm-ingress-router/types"
	"github.com/tpbowden/swarm-ingress-router/version"
)

type CLI struct {
	newServer    func(string, string, string, string, int, int) types.Startable
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
        cli.StringFlag{
          Name:  "cert, c",
          Value: "",
          Usage: "Global wildcard certificate for services.",
        },
        cli.StringFlag{
          Name:  "key, k",
          Value: "",
          Usage: "Global key for wildcard cert.",
        },
				cli.IntFlag{
					Name:  "max-body-size",
					Value: 4,
					Usage: "Max body size in MB",
				},
        cli.IntFlag{
          Name:  "read-buffer-size",
          Value: 4,
          Usage: "Per-connection size in KB",
        },
			},
			Action: func(ctx *cli.Context) error {
        // Default read buffer size is currently 4096 (defaultReadBufferSize in
        // fasthttp's server.go)
        // Per-connection buffer size for requests' reading.
        // This also limits the maximum header size.
        //
        // Increase this buffer if your clients send multi-KB RequestURIs
        // and/or multi-KB headers (for example, BIG cookies, things like kerberos
        // and others which use larger sizes than 4KB)
				server := c.newServer(ctx.String("bind"), ctx.GlobalString("redis"), ctx.String("cert"), ctx.String("key"), ctx.Int("max-body-size")*1024*1024, ctx.Int("read-buffer-size")*1024)
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
