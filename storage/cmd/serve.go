package cmd

import (
	"log"

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
	s, err := core.NewServerGRPC(core.ServerGRPCConfig{
		Port: c.Int("port"),
	})
	if err != nil {
		log.Fatalf("Error to creating server: %v", err)
	}

	if err := s.Listen(); err != nil {
		log.Fatalf("Error on server listen: %v", err)
	}

	defer s.Close()

	return nil
}
