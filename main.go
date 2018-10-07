package main

import (
	"os"

	"github.com/evanharmon/eph-music-micro/storage/cmd"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:  "eph-music",
		Usage: "use that music",
		Commands: []*cli.Command{
			&cmd.Serve,
			&cmd.Upload,
		},
	}

	app.Run(os.Args)
}
