// +build integration

// Package main for Cloud Storage api interactions
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	gstorage "cloud.google.com/go/storage"
	helper "github.com/evanharmon/eph-music-micro/helper"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

var (
	envs = map[string]string{
		"credentials": "GOOGLE_APPLICATION_CREDENTIALS",
		"ProjectID":   "GOOGLE_PROJECT_ID"}
	credentials string
)

type Service interface {
	Create(ctx context.Context, p *pb.Project) error
	Delete(ctx context.Context) error
	ListBuckets(p *pb.Project) ([]string, error)
	UploadFile(ctx context.Context, path string) error
	DeleteFile(ctx context.Context, name string) error
}

type StorageBucket struct {
	client *gstorage.Client
	handle *gstorage.BucketHandle
	// bucketName is coupled to the handle
	name string
}

// Init function loads required environment variables
func Init() {
	for k, v := range envs {
		export, err := helper.GetEnv(v)
		if err != nil {
			log.Fatal(err)
		}
		envs[k] = export
	}
}

// ConfigureStorage creates a client for re-use.
// The client is not tied to a project id.
func configure() (*gstorage.Client, error) {
	client, err := gstorage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// New inits and returns the bucket handler and client
func New(projectID string, name string) (Service, error) {
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
	return &StorageBucket{client, handle, name}, nil
}

// ListBuckets provides a way to list all storage buckets by Project ID.
// Change the `ProjectID` package global for other project bucket lists
func (s *StorageBucket) ListBuckets(p *pb.Project) ([]string, error) {
	var buckets []string
	it := s.client.Buckets(context.Background(), p.Id)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Bucket iterator failed: %v", err)
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets, nil
}

// Create attempts to create a new bucket for a project id
// this custom create function is idempotent and will not return the 409 error
// if bucket is owned and already exists
func (s *StorageBucket) Create(ctx context.Context, p *pb.Project) error {
	err := s.handle.Create(ctx, p.Id, nil)
	gerr, ok := err.(*googleapi.Error)
	if err != nil && !ok {
		return err
	}

	msg := "You already own this bucket. Please select another name."
	if err != nil && ok && gerr.Message != msg {
		return err
	}

	return nil
}

// Delete the bucket
func (s *StorageBucket) Delete(ctx context.Context) error {
	if err := s.handle.Delete(ctx); err != nil {
		return fmt.Errorf("Failed to delete storage bucket %s: %v", s.name, err)
	}

	return nil
}

// UploadFile to storage bucket
func (s *StorageBucket) UploadFile(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(f)

	fname := filepath.Base(path)
	wc := s.handle.Object(fname).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

// DeleteFile from storage bucket
func (s *StorageBucket) DeleteFile(ctx context.Context, name string) error {
	if err := s.handle.Object(name).Delete(ctx); err != nil {
		return err
	}
	return nil
}
