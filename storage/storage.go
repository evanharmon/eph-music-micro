// Package storage for Cloud Storage api interactions
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"cloud.google.com/go/storage"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

// ProjectID can be changed to re-use the client
var (
	envs = map[string]string{
		"credentials": "GOOGLE_APPLICATION_CREDENTIALS",
		"projectID":   "GOOGLE_PROJECT_ID"}
	credentials       string
	ProjectID         string
	storageBucket     *storage.BucketHandle
	storageBucketName string
	storageClient     *storage.Client
)

// Init function loads required environment variables
func Init() {
	for k, v := range envs {
		v, err := getEnv(v)
		if err != nil {
			log.Fatal(err)
		}
		envs[k] = v
	}
}

// getEnv function provides a safe lookup for environment variables
func getEnv(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("Env variable must be provided to getEnv")
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Could not find environment variable: %s", key)
	}

	return val, nil
}

func checkConfig() error {
	if storageBucket == nil {
		return errors.New("Use ConfigureStore() before calling ListBuckets()")
	}
	if ProjectID == "" {
		return errors.New("ProjectID must not be an empty string")
	}

	return nil
}

// ConfigureStorage creats a client re-use.
// The client is not tied to a project id.
func ConfigureStorage(bucketName string) error {
	if len(bucketName) == 0 {
		return errors.New("BucketName must be provided")
	}
	storageBucketName = bucketName

	ctx := context.Background()
	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		return err
	}

	storageBucket = storageClient.Bucket(storageBucketName)
	return nil
}

// ListBuckets provides a way to list all storage buckets by Project ID.
// Change the `ProjectID` package global for other project bucket lists
func ListBuckets() ([]string, error) {
	if err := checkConfig(); err != nil {
		return []string{}, err
	}
	ctx := context.Background()
	var buckets []string
	it := storageClient.Buckets(ctx, ProjectID)
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
// if bucket already exists
func Create(ctx context.Context) error {
	if err := checkConfig(); err != nil {
		return err
	}
	err := storageBucket.Create(ctx, ProjectID, nil)
	gerr, ok := err.(*googleapi.Error)
	if err != nil && ok && gerr.Code != 409 {
		return err
	}

	return nil
}

// Delete the bucket
func Delete(ctx context.Context) error {
	if err := checkConfig(); err != nil {
		return err
	}
	if err := storageBucket.Delete(ctx); err != nil {
		return fmt.Errorf("Failed to delete storage bucket %s: %v", storageBucketName, err)
	}

	return nil
}

// UploadFile to storage bucket
func UploadFile(ctx context.Context, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fname)
	wc := storageBucket.Object(name).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
