// Package storage for Cloud Storage api interactions
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	helper "github.com/evanharmon/eph-music-micro/helper"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

var (
	envs = map[string]string{
		"credentials": "GOOGLE_APPLICATION_CREDENTIALS",
		"ProjectID":   "GOOGLE_PROJECT_ID"}
	credentials string
)

type StorageBucket struct {
	handle    *storage.BucketHandle
	name      string
	client    *storage.Client
	projectID string
}

// New inits and returns the bucket handler and client
func New(projectID string, name string) (*StorageBucket, error) {
	if projectID == "" {
		return nil, errors.New("ProjectID must not be an empty string")
	}

	if len(name) == 0 {
		return nil, errors.New("BucketName must be provided")
	}

	client, err := configure(name)
	if err != nil {
		return nil, err
	}
	handle := client.Bucket(name)
	return &StorageBucket{handle, name, client, projectID}, nil
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
func configure(name string) (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ListBuckets provides a way to list all storage buckets by Project ID.
// Change the `ProjectID` package global for other project bucket lists
func (bkt *StorageBucket) ListBuckets() ([]string, error) {
	ctx := context.Background()
	var buckets []string
	it := bkt.client.Buckets(ctx, bkt.projectID)
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
func (bkt *StorageBucket) Create(ctx context.Context) error {
	err := bkt.handle.Create(ctx, bkt.projectID, nil)
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
func (bkt *StorageBucket) Delete(ctx context.Context) error {
	if err := bkt.handle.Delete(ctx); err != nil {
		return fmt.Errorf("Failed to delete storage bucket %s: %v", bkt.name, err)
	}

	return nil
}

// UploadFile to storage bucket
func (bkt *StorageBucket) UploadFile(ctx context.Context, fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(f)

	fname := filepath.Base(fpath)
	wc := bkt.handle.Object(fname).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

// DeleteFile from storage bucket
func (bkt *StorageBucket) DeleteFile(ctx context.Context, name string) error {
	if err := bkt.handle.Object(name).Delete(ctx); err != nil {
		return err
	}
	return nil
}
