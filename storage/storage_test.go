package storage

import (
	"context"
	"os"
	"testing"

	test "github.com/evanharmon/eph-music-micro/storage/testhelper"
)

var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	projectID       = "eph-music"
	bucketName      = "eph-music"
)

// BeforeTestGetEnv unsets the environment variable for a Google Cloud Project
// ID
func BeforeTestGetEnv(t *testing.T) {
	test.Ok(t, os.Unsetenv(envVarProjectID))
}

// TestGetEnv tests the safe lookup of environment variables
func TestGetEnv(t *testing.T) {
	// Test Empty String
	BeforeTestGetEnv(t)
	val, err := getEnv("")
	test.Assert(t, len(val) == 0, "Empty string as argument should return empty string as default")
	test.Throws(t, err)

	// Test ENV missing
	BeforeTestGetEnv(t)
	val2, err2 := getEnv(envVarProjectID)
	test.Assert(t, len(val2) == 0, "Missing ENV var should return empty string as default")
	test.Throws(t, err2)

	// VALID ENV should return
	BeforeTestGetEnv(t)
	os.Setenv(envVarProjectID, "")
	val4, err4 := getEnv(envVarProjectID)
	test.Assert(t, len(val4) == 0, "Valid ENV set should get value")
	test.Ok(t, err4)
}

func BeforeConfigureStorage(t *testing.T) {
	ProjectID = projectID
	return
}

// TestConfigureStoraage tests creating a client for the Google Cloud Storage
// API
func TestConfigureStorage(t *testing.T) {
	BeforeConfigureStorage(t)
	test.Ok(t, ConfigureStorage(bucketName))
}

func BeforeTestCreate(t *testing.T) {
	ProjectID = projectID
	test.Ok(t, ConfigureStorage(bucketName))
}

func TestCreate(t *testing.T) {
	// test initial bucket create
	BeforeTestCreate(t)
	ctx := context.Background()
	test.Ok(t, Create(ctx))

	// test duplicate bucket creation
	test.Ok(t, Create(ctx))
}

func BeforeTestDelete(t *testing.T) {
	ProjectID = projectID
	test.Ok(t, ConfigureStorage(bucketName))

	ctx := context.Background()
	test.Ok(t, Create(ctx))
}

func TestDelete(t *testing.T) {
	BeforeTestDelete(t)
	ctx := context.Background()
	test.Ok(t, Delete(ctx))
}

// func TestUploadFile(t *testing.T) {

// }
