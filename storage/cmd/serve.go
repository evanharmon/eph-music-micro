package cmd

import (
	"fmt"

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
	},
}

func serveAction(c *cli.Context) error {
	s, err := core.NewProviderGRPC(core.ProviderGRPCConfig{
		Port: c.Int("port"),
	})
	if err != nil {
		fmt.Printf("Error creating server: %v", err)
		return cli.Exit(err, 1)
	}

	if err := s.Listen(); err != nil {
		fmt.Printf("Error on server listen: %v", err)
		return cli.Exit(err, 1)
	}

	defer s.Close()

	return nil
}
