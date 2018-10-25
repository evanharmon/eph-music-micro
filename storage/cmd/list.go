// Package cmd provides ...
package cmd

import (
	"context"
	"errors"

	"github.com/evanharmon/eph-music-micro/storage/core"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	cli "gopkg.in/urfave/cli.v2"
)

var ListBuckets = cli.Command{
	Name:   "listbuckets",
	Usage:  "list buckets",
	Action: listAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "project",
			Usage: "project id",
			Value: "evan-terraform-admin",
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "address",
			Value: "localhost:10013",
		},
	},
}

func listAction(c *cli.Context) error {
	var (
		err error

		address = c.String("address")
		client  = core.ClientGRPC{}
		project = c.String("project")
	)

	if address == "" {
		err = errors.New("Address is required")
		return cli.Exit(err, 1)
	}

	if project == "" {
		err = errors.New("project is required")
		return cli.Exit(err, 1)
	}

	client, err = core.NewClientGRPC(core.ClientGRPCConfig{
		Address: address,
	})
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer client.Close()

	_, err = client.ListBuckets(context.Background(), &pb.ListBucketsRequest{
		Project: &pb.Project{Id: project},
	})
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}
