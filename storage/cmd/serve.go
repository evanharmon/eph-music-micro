package cmd

import (
	"github.com/evanharmon/eph-music-micro/storage/core"
	cli "gopkg.in/urfave/cli.v2"
)

var Serve = cli.Command{
	Name:   "serve",
	Usage:  "initiates a gRPC server",
	Action: serveAction,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Usage: "port to bind to",
			Value: 10013,
		},
		&cli.StringFlag{
			Name:  "project",
			Usage: "google project",
			Value: "test-eph-music",
		},
		&cli.StringFlag{
			Name:  "bucket",
			Usage: "bucket name",
			Value: "test-eph-music",
		},
	},
}

func serveAction(c *cli.Context) error {
	var (
		err error

		server  = &core.ServerGRPC{}
		port    = c.Int("port")
		project = c.String("project")
		name    = c.String("bucket")
	)

	server, err = core.NewServerGRPC(core.ServerGRPCConfig{
		Port:    port,
		Project: project,
		Name:    name,
	})
	must(err)
	server.Listen()
	must(err)
	defer server.Close()

	return nil
}
