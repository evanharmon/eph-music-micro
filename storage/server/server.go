package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	gstorage "cloud.google.com/go/storage"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
)

const port = 50051

type ServerGRPC struct {
	server *grpc.Server
	client *gstorage.Client
	handle *gstorage.BucketHandle
	// bucketName is coupled to the handle
	name string
}

// ListBuckets provides a way to list all storage buckets by Project ID.
// Change the `ProjectID` package global for other project bucket lists
func (s *ServerGRPC) ListBuckets(ctx context.Context, req *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error) {
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
	res := &pb.ListBucketsResponse{Buckets: buckets}
	return res, nil
}

func (s *ServerGRPC) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	handle := s.client.Bucket(req.Bucket.Name)

	err := handle.Create(ctx, req.Project.Id, nil)
	gerr, ok := err.(*googleapi.Error)
	if err != nil && !ok {
		return nil, err
	}

	msg := "You already own this bucket. Please select another name."
	if err != nil && ok && gerr.Message != msg {
		fmt.Println("Bucket already created")
		return nil, err
	}

	res := &pb.CreateResponse{Result: "success"}
	return res, nil
}

// Delete the bucket
func (s *ServerGRPC) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	handle := s.client.Bucket(req.Bucket.Name)
	if err := handle.Delete(ctx); err != nil {
		return nil, fmt.Errorf("Failed to delete storage bucket %s: %v", s.name, err)
	}
	res := &pb.DeleteResponse{Result: "success"}
	return res, nil
}

// UploadFile to storage bucket
func (s *ServerGRPC) UploadFile(stream pb.Storage_UploadFileServer) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			return err
		}
	}

END:
	err := stream.SendAndClose(&pb.UploadFileResponse{
		Message: "Upload received with success",
		Code:    pb.UploadStatusCode_Ok,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile from storage bucket
func (s *ServerGRPC) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	handle := s.client.Bucket(req.Bucket.Name)
	if err := handle.Object(req.File.Name).Delete(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}

func configure() (*gstorage.Client, error) {
	client, err := gstorage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// New inits and returns the bucket handler and client
func NewServerGRPC(projectID string, name string) (*ServerGRPC, error) {
	if projectID == "" {
		return nil, errors.New("ProjectID must not be an empty string")
	}

	if len(name) == 0 {
		return nil, errors.New("BucketName must be provided")
	}

	client, err := configure()
	if err != nil {
		return nil, err
	}
	handle := client.Bucket(name)
	server := grpc.NewServer()
	s := &ServerGRPC{server, client, handle, name}
	pb.RegisterStorageServer(server, s)

	return s, nil
}

func (s *ServerGRPC) Listen() error {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("Failed to serve: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("starting app...")
	fmt.Printf("Listening on port: %v", port)

	s, err := NewServerGRPC("evan-terraform-admin", "evan-terraform-admin")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Listen(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
