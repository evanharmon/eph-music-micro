package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type ClientGRPC struct {
	conn      *grpc.ClientConn
	client    pb.StorageClient
	chunkSize int
}

// ListBuckets provides a way to list all storage buckets by Project ID.
// Change the `ProjectID` package global for other project bucket lists
func (c *ClientGRPC) runListBuckets(ctx context.Context, req *pb.ListBucketsRequest) error {
	res, err := c.client.ListBuckets(ctx, req)
	if err != nil {
		return errors.Errorf("%v.GetBuckets(_) = _, %v", c.client, err)
	}
	log.Printf("Response from ListBuckets: %v", res.Buckets)

	return nil
}

// Create attempts to create a new bucket for a project id
// this custom create function is idempotent and will not return the 409 error
// if bucket is owned and already exists
func (c *ClientGRPC) runCreate(ctx context.Context, req *pb.CreateRequest) error {
	res, err := c.client.Create(ctx, req)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

// Delete the bucket
func (c *ClientGRPC) runDelete(ctx context.Context, req *pb.DeleteRequest) error {
	res, err := c.client.Delete(ctx, req)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

// UploadFile to storage bucket
func (c *ClientGRPC) runUploadFile(ctx context.Context, fname string) error {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *pb.UploadFileResponse
	)

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		return err
	}
	defer func(s pb.Storage_UploadFileClient) {
		err = stream.CloseSend()
		if err != nil {
			log.Fatal(err)
		}
	}(stream)

	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err != io.EOF {
				writing = false
				err = nil
				continue
			}

			return nil
		}

		err = stream.Send(&pb.UploadFileRequest{
			Project: &pb.Project{
				Id: "evan-terraform-admin",
			},
			Bucket: &pb.Bucket{
				Name: "eph-test-music",
			},
			Chunk: &pb.Chunk{
				Content: buf[:n],
			},
		})
		if err != nil {
			return err
		}
	}

	status, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	if status.Code != pb.UploadStatusCode_Ok {
		return fmt.Errorf("upload failed - msg: %s", status.Message)
	}
	return nil
}

// DeleteFile from storage bucket
func (c *ClientGRPC) runDeleteFile(ctx context.Context, req *pb.DeleteFileRequest) error {
	return nil
}

func (c *ClientGRPC) Close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

func NewClientGRPC() (c ClientGRPC, err error) {
	c.chunkSize = 1024
	c.conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return ClientGRPC{}, fmt.Errorf("Could not connect: %v", err)
	}
	c.client = pb.NewStorageClient(c.conn)

	return c, nil
}

func main() {
	c, err := NewClientGRPC()
	if err != nil {
		log.Fatal(err)
	}

	// lbreq := &pb.ListBucketsRequest{
	// Project: &pb.Project{
	// Id: "evan-terraform-admin",
	// },
	// }
	// err = c.runListBuckets(context.Background(), lbreq)
	// if err != nil {
	// log.Fatal(err)
	// }

	creq := &pb.CreateRequest{
		Project: &pb.Project{
			Id: "evan-terraform-admin",
		},
		Bucket: &pb.Bucket{
			Name: "test-eph-music",
		},
	}
	err = c.runCreate(context.Background(), creq)
	if err != nil {
		log.Fatal(err)
	}

	// fname := "/Users/evan/go/src/github.com/evanharmon/eph-music-micro/storage/testdata/upload-file.txt"
	// err = c.runUploadFile(context.Background(), fname)
	// if err != nil {
	// log.Fatal(err)
	// }

	dreq := &pb.DeleteRequest{
		Project: &pb.Project{
			Id: "evan-terraform-admin",
		},
		Bucket: &pb.Bucket{
			Name: "test-eph-music",
		},
	}
	err = c.runDelete(context.Background(), dreq)
	if err != nil {
		log.Fatal(err)
	}
}
