package core

//go:generate mockgen -destination mocks/mock_server.go -package mocks github.com/evanharmon/eph-music-micro/storage/core ServerService

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	gstorage "cloud.google.com/go/storage"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
)

type ProviderService interface {
	NewProviderGRPC(cfg ProviderGRPCConfig) (*ProviderGRPC, error)
	Listen(port int) (net.Listener, error)
	Server(lis net.Listener) error
	Close()
	ListBuckets(context.Context, *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error)
	Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error)
	Delete(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error)
	UploadFile(*pb.Storage_UploadFileServer) error
	DeleteFile(context.Context, *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error)
}

type ProviderGRPC struct {
	client *gstorage.Client
	server *grpc.Server
	port   int
}

type ProviderGRPCConfig struct {
	Port int
}

// NewProviderGRPC creates a new grpc server
func NewProviderGRPC(cfg ProviderGRPCConfig) (*ProviderGRPC, error) {
	var (
		port = cfg.Port
	)
	if port == 0 {
		return nil, errors.New("Port must be specified")
	}

	client, err := gstorage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	s := &ProviderGRPC{client, server, port}
	pb.RegisterStorageServer(server, s)

	return s, nil
}

func (s *ProviderGRPC) Listen() error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		s.Close()
		return fmt.Errorf("Failed to listen: %v", err)
	}
	fmt.Printf("Server listening on port: %v\n", s.port)

	if err := s.server.Serve(lis); err != nil {
		s.Close()
		return fmt.Errorf("Failed to serve: %v", err)
	}
	return nil
}

func (s *ProviderGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}
	return
}

// ListBuckets provides a way to list all storage buckets by Project ID.
func (s *ProviderGRPC) ListBuckets(ctx context.Context, req *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error) {
	if req.Project.Id == "" {
		return nil, errors.New("Project ID is required")
	}
	var buckets []*pb.Bucket
	it := s.client.Buckets(ctx, req.Project.Id)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket iterator failed: %v", err)
		}
		buckets = append(buckets, &pb.Bucket{
			Name: battrs.Name,
		})
	}
	return &pb.ListBucketsResponse{Buckets: buckets}, nil
}

// Create the bucket
func (s *ProviderGRPC) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	bkt := s.client.Bucket(req.Bucket.Name)
	err := bkt.Create(ctx, req.Project.Id, nil)
	gerr, ok := err.(*googleapi.Error)
	if err != nil && !ok {
		return nil, err
	}

	msg := "You already own this bucket. Please select another name."
	if err != nil && ok && gerr.Message != msg {
		fmt.Println("Bucket already created")
		return nil, err
	}

	return &pb.CreateResponse{Result: "success"}, nil
}

// Delete the bucket
func (s *ProviderGRPC) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	bkt := s.client.Bucket(req.Bucket.Name)
	if err := bkt.Delete(ctx); err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Result: "success"}, nil
}

// UploadFile to storage bucket
// Request Protobuf only available via the stream
// Response Protobuf is sent back via the closing of the stream
func (s *ProviderGRPC) UploadFile(stream pb.Storage_UploadFileServer) error {
	var (
		buf        []byte
		err        error
		fileName   string
		bucketName string
	)
	for {
		// BEWARE last iteration of Recv(): req = nil, err = io.EOF
		req, err := stream.Recv()
		// only set fname once
		if req != nil && fileName == "" {
			fileName = req.File.Name
		}

		if req != nil && bucketName == "" {
			bucketName = req.Bucket.Name
		}

		if err != nil {
			if err == io.EOF {
				err = nil
				goto END
			}

			return err
		}

		buf = append(buf, req.Chunk.Content...)
	}

END:
	// Implement io.Reader on buf and copy to object
	nr := bytes.NewReader(buf)
	bkt := s.client.Bucket(bucketName)
	wc := bkt.Object(fileName).NewWriter(context.Background())
	if _, err = io.Copy(wc, nr); err != nil {
		err = stream.SendAndClose(&pb.UploadFileResponse{
			Message: "Upload failed to copy NewReader to NewWriter",
			Code:    pb.UploadStatusCode_Failed,
		})
		if err != nil {
			return err
		}
	}
	// Close and Upload
	if err = wc.Close(); err != nil {
		err = stream.SendAndClose(&pb.UploadFileResponse{
			Message: "Upload failed closing NewWriter",
			Code:    pb.UploadStatusCode_Failed,
		})
		if err != nil {
			return err
		}
	}
	err = stream.SendAndClose(&pb.UploadFileResponse{
		Message: "Upload received with success",
		Code:    pb.UploadStatusCode_Ok,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile from storage bucket
func (s *ProviderGRPC) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	if req.File.Name == "" {
		return nil, fmt.Errorf("File name to delete cannot be an empty string")
	}

	bkt := s.client.Bucket(req.Bucket.Name)
	if err := bkt.Object(req.File.Name).Delete(ctx); err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{Result: "success"}, nil
}
