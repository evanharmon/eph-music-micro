package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
)

var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	projectID       = "eph-music"
	bucketName      = "eph-music"
)

// BeforeTestGetEnv unsets the environment variable for a Google Cloud Project
// ID
func BeforeTestGetEnv() error {
	if err := os.Unsetenv(envVarProjectID); err != nil {
		return fmt.Errorf("BeforeTestGetEnv: failed to unset ENV %s: %v", envVarProjectID, err)
	}

	return nil
}

// TestGetEnv tests the safe lookup of environment variables
func TestGetEnv(t *testing.T) {
	// Test Empty String
	if err := BeforeTestGetEnv(); err != nil {
		t.Error(err)
	}

	val, err := getEnv("")
	if len(val) != 0 {
		t.Error("Empty string as argument should return empty string as default")
	}
	if err == nil {
		t.Error("Empty string as argument string should throw error")
	}

	// Test ENV missing
	if err := BeforeTestGetEnv(); err != nil {
		t.Error(err)
	}

	val2, err2 := getEnv(envVarProjectID)
	if len(val2) != 0 {
		t.Error("Missing ENV var should return empty string as default")
	}
	if err2 == nil {
		t.Error("Missing ENV var should throw error")
	}

	// VALID ENV should return
	if err := BeforeTestGetEnv(); err != nil {
		t.Error(err)
	}

	os.Setenv(envVarProjectID, "")
	val4, err4 := getEnv(envVarProjectID)
	if err4 != nil {
		t.Error("Valid ENV set should not throw error")
	}
	if len(val4) != 0 {
		t.Error("Valid ENV set should get value")
	}
}

func BeforeConfigureStorage() error {
	ProjectID = projectID
	return nil
}

// TestConfigureStoraage tests creating a client for the Google Cloud Storage
// API
func TestConfigureStorage(t *testing.T) {
	if err := BeforeConfigureStorage(); err != nil {
		t.Error(err)
	}

	err := ConfigureStorage(bucketName)
	if err != nil {
		t.Errorf("should not throw error on valid bucketId %s\n", bucketName)
	}
}

func BeforeTestCreate() error {
	ProjectID = projectID
	if err := ConfigureStorage(bucketName); err != nil {
		return fmt.Errorf("BeforeTestCreate failed to configureStorage: %v", err)
	}

	return nil
}

func TestCreate(t *testing.T) {
	// test initial bucket create
	if err := BeforeTestCreate(); err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	if err := Create(ctx); err != nil {
		t.Errorf("Failed to create new bucket: %v", err)
	}

	// test duplicate bucket creation
	if err := Create(ctx); err != nil {
		t.Errorf("Failed to create new bucket: %v", err)
	}
}

func BeforeTestDelete() error {
	ProjectID = projectID
	if err := ConfigureStorage(bucketName); err != nil {
		return fmt.Errorf("BeforeTestCreate failed to configureStorage: %v", err)
	}

	ctx := context.Background()
	if err := Create(ctx); err != nil {
		return fmt.Errorf("BeforeTestDelete failed to create new bucket: %v", err)
	}

	return nil
}

func TestDelete(t *testing.T) {
	if err := BeforeTestDelete(); err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	if err := Delete(ctx); err != nil {
		t.Errorf("Failed to delete bucket: %v", err)
	}
}
