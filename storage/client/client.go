package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

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
func (c *ClientGRPC) runUploadFile(ctx context.Context, req *pb.UploadFileRequest) error {
	var (
		buf     []byte
		err     error
		file    *os.File
		n       int
		status  *pb.UploadFileResponse
		writing = true
	)

	file, err = os.Open(req.File.Path)
	if err != nil {
		return fmt.Errorf("Error opening file: %v\n", err)
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		return fmt.Errorf("Error uploading file via stream: %v\n", stream)
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
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			return fmt.Errorf("Error copying from file to buf: %v\n", err)
		}

		req.Chunk.Content = buf[:n]
		err = stream.Send(req)
		if err != nil {
			return fmt.Errorf("Error on stream.Send() %v\n", err)
		}
	}

	status, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("Failed to receive upstream status response: %v\n", err)
	}

	if status.Code != pb.UploadStatusCode_Ok {
		return fmt.Errorf("upload failed - msg: %s\n", status.Message)
	}

	fmt.Printf("upload file - msg: %s\n", status.Message)
	return nil
}

// DeleteFile from storage bucket
func (c *ClientGRPC) runDeleteFile(ctx context.Context, req *pb.DeleteFileRequest) error {
	res, err := c.client.DeleteFile(ctx, req)
	if err != nil {
		return err
	}
	log.Println(res)

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
	var (
		fpath      string
		fname      string
		err        error
		projectId  string
		bucketName string
	)
	fpath = "/Users/evan/go/src/github.com/evanharmon/eph-music-micro/storage/testdata/upload-file.txt"
	fname = path.Base(fpath)
	projectId = "evan-terraform-admin"
	bucketName = "test-eph-music"

	c, err := NewClientGRPC()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listing Buckets")
	err = c.runListBuckets(context.Background(), &pb.ListBucketsRequest{
		Project: &pb.Project{Id: projectId},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating Bucket")
	err = c.runCreate(context.Background(), &pb.CreateRequest{
		Project: &pb.Project{Id: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Uploading File")
	err = c.runUploadFile(context.Background(), &pb.UploadFileRequest{
		Project: &pb.Project{Id: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		File:    &pb.File{Name: fname, Path: fpath},
		Chunk:   &pb.Chunk{Content: []byte{}},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Delete File")
	err = c.runDeleteFile(context.Background(), &pb.DeleteFileRequest{
		Project: &pb.Project{Id: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		File:    &pb.File{Name: fname, Path: fpath},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Delete Bucket")
	err = c.runDelete(context.Background(), &pb.DeleteRequest{
		Project: &pb.Project{Id: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	})
	if err != nil {
		log.Fatal(err)
	}
}
