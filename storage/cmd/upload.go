package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/evanharmon/eph-music-micro/storage/core"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	cli "gopkg.in/urfave/cli.v2"
)

var Upload = cli.Command{
	Name:   "upload",
	Usage:  "upload a file to a storage bucket",
	Action: uploadAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Usage: "file name",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "address of the server to connect to",
			Value: "localhost:10013",
		},
		&cli.IntFlag{
			Name:  "chunk-size",
			Usage: "size of the chunk messages (grpc only)",
			Value: (1 << 12),
		},
		&cli.StringFlag{
			Name:  "project",
			Usage: "project id",
			Value: "eph-music",
		},
		&cli.StringFlag{
			Name:  "bucket",
			Usage: "bucket name",
			Value: "test-eph-music",
		},
	},
}

func uploadAction(c *cli.Context) error {
	var (
		err   error
		fpath string
		fname string

		address   = c.String("address")
		chunkSize = c.Int("chunk-size")
		client    = core.ClientGRPC{}
		file      = c.String("file")
		project   = c.String("project")
		bucket    = c.String("bucket")
	)

	if address == "" {
		errors.New("Address is required")
		return cli.Exit(err, 1)
	}

	if file == "" {
		errors.New("file must be set")
		return cli.Exit(err, 1)
	}
	fpath, err = filepath.Abs(file)
	if err != nil {
		return cli.Exit(fmt.Errorf("File not found: %s", file), 1)
	}
	fname = filepath.Base(file)

	client, err = core.NewClientGRPC(core.ClientGRPCConfig{
		Address:   address,
		ChunkSize: chunkSize,
	})
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer client.Close()

	_, err = client.UploadFile(context.Background(), &pb.UploadFileRequest{
		Project: &pb.Project{Id: project},
		Bucket:  &pb.Bucket{Name: bucket},
		File:    &pb.File{Name: fname, Path: fpath},
		Chunk:   &pb.Chunk{Content: []byte{}},
	})
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}
