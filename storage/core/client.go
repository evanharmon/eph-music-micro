package core

//go:generate mockgen -destination mocks/mock_client.go -package mocks github.com/evanharmon/eph-music-micro/storage/core ClientService

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

type ClientService interface {
	NewClientGRPC(ClientGRPCConfig) (ClientGRPC, error)
	Close()
	ListBuckets(context.Context, *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error)
	Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error)
	Delete(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error)
	UploadFile(context.Context, *pb.UploadFileRequest) (*pb.UploadFileResponse, error)
	DeleteFile(context.Context, *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error)
}

type ClientGRPC struct {
	conn      *grpc.ClientConn
	client    pb.StorageClient
	chunkSize int
}

type ClientGRPCConfig struct {
	Address   string
	ChunkSize int
}

func NewClientGRPC(cfg ClientGRPCConfig) (ClientGRPC, error) {
	var (
		err error
		c   ClientGRPC

		grpcOpts = []grpc.DialOption{}
	)
	// Certs Not Implemented
	grpcOpts = append(grpcOpts, grpc.WithInsecure())

	if cfg.Address == "" {
		return c, errors.Errorf("address must be specified")
	}

	// Cleaner Than IF statement
	switch {
	case cfg.ChunkSize == 0:
		return c, errors.Errorf("Chunksize must be specified")
	case cfg.ChunkSize > (1 << 22):
		return c, errors.Errorf("Chunksize must be less than 4MB")
	default:
		c.chunkSize = cfg.ChunkSize
	}

	c.conn, err = grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		return c, errors.Wrapf(err, "Failed to start grpc connection with address: %s", cfg.Address)
	}

	c.client = pb.NewStorageClient(c.conn)

	return c, nil
}

func (c *ClientGRPC) Close() {
	if c.conn == nil {
		return
	}
	if err := c.conn.Close(); err != nil {
		log.Fatal(err)
	}
}

// ListBuckets provides a way to list all storage buckets by Project ID.
func (c *ClientGRPC) ListBuckets(ctx context.Context, req *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error) {
	res, err := c.client.ListBuckets(ctx, req)
	if err != nil {
		return nil, errors.Errorf("%v.GetBuckets(_) = _, %v", c.client, err)
	}
	log.Printf("Response from ListBuckets: %v", res.Buckets)

	return res, nil
}

// Create attempts to create a new bucket for a project id
// this custom create function is idempotent and will not return the 409 error
// if bucket is owned and already exists
func (c *ClientGRPC) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	res, err := c.client.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Println(res)

	return res, nil
}

// Delete the bucket
func (c *ClientGRPC) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	res, err := c.client.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Println(res)

	return res, nil
}

// UploadFile to storage bucket
func (c *ClientGRPC) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	var (
		buf     []byte
		err     error
		n       int
		status  *pb.UploadFileResponse
		writing = true
	)

	file, err := os.Open(req.File.Path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v\n", err)
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error uploading file via stream: %v\n", stream)
	}
	defer func(s pb.Storage_UploadFileClient) {
		if err = s.CloseSend(); err != nil {
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
			return nil, fmt.Errorf("Error copying from file to buf: %v\n", err)
		}

		req.Chunk.Content = buf[:n]
		err = stream.Send(req)
		if err != nil {
			return nil, fmt.Errorf("Error on stream.Send() %v\n", err)
		}
	}

	status, err = stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("Failed to receive upstream status response: %v\n", err)
	}

	if status.Code != pb.UploadStatusCode_Ok {
		return nil, fmt.Errorf("upload failed - msg: %s\n", status.Message)
	}

	return status, nil
}

// DeleteFile from storage bucket
func (c *ClientGRPC) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	res, err := c.client.DeleteFile(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
