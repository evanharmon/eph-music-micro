// Package storage for Cloud Storage api interactions
package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	envs              = map[string]string{"credentials": "GOOGLE_APPLICATION_CREDENTIALS", "projectID": "GOOGLE_PROJECT_ID"}
	credentials       string
	projectID         string
	storageBucket     *storage.BucketHandle
	storageBucketName string
	// Client should be re-used and is safe for concurrent use
	storageClient *storage.Client
)

// init function loads required environment variables
func init() {
	for k, v := range envs {
		v, err := getEnv(v)
		if err != nil {
			log.Fatalf("Environment variable %s must be set", k)
		}
		envs[k] = v
	}
}

// getEnv function provides a safe lookup for environment variables
func getEnv(key string) (string, error) {
	fmt.Printf("looking up env: %s\n", key)
	if len(key) == 0 {
		return "", errors.New("Key must be provided to getEnv")
	}
	val, ok := os.LookupEnv(key)
	fmt.Printf("LookupEnv val: %v\n", val)
	fmt.Printf("LookupEnv ok: %v\n", ok)
	if !ok {
		return "", fmt.Errorf("Could not find environment variable: %s", key)
	}

	return val, nil
}

// configureStorage provides a re-usable client for the Google Cloud Storage API
func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	if len(bucketID) == 0 {
		return nil, errors.New("bucketID must be provided")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.Bucket(bucketID), nil
}

// listBucketsByProjectID provides a way to list all storage buckets in a Google
// Cloud Project ID
func listBucketsByProjectID(projectID string) ([]string, error) {
	if projectID == "" {
		return []string{}, errors.New("No valid projectID found")
	}
	ctx := context.Background()
	var buckets []string
	it := storageClient.Buckets(ctx, projectID)
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
