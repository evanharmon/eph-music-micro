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

const (
	// TODO ADD GOOGLE_APPLICATION_CREDENTIALS
	envVarProjectID = "GOOGLE_PROJECT_ID" // Follow Google's gcloud cli Export naming convention
)

var (
	projectID         string
	storageBucket     *storage.BucketHandle
	storageBucketName string
	// Client should be re-used and is safe for concurrent use
	storageClient *storage.Client
)

// TODO ADD GOOGLE_APPLICATION_CREDENTIALS
// maybe in an init function? since when binary launches those env vars should
// exist anyways

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

// getProjectID provides a safe lookup for the Google Cloud Project ID
func getProjectID() (string, error) {
	val, err := getEnv(envVarProjectID)
	if err != nil {
		return "", err
	}
	if len(val) == 0 {
		return "", errors.New("GOOGLE_PROJECT_ID should not be an empty string")
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
	projectID, err := getProjectID()
	if err != nil {
		return []string{}, err
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
