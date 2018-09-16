package storage

import (
	"context"
	"fmt"
	"testing"

	helper "github.com/evanharmon/eph-music-micro/helper"
	test "github.com/evanharmon/eph-music-micro/helper/testhelper"
)

//
var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	bucketName      = ""
)

// Setup sets the necessary vars from ENV exports
func setup(t *testing.T) {
	id, err := helper.GetEnv(envVarProjectID)
	test.Ok(t, err)
	test.Assert(t, len(id) != 0, "Export necessary ENV vars before running tests")
	ProjectID = id
	bucketName = fmt.Sprintf("%s-test", id)
}

func beforeConfigureStorage(t *testing.T) {
	setup(t)
}

// TestConfigureStoraage tests creating a client for the Google Cloud Storage
// API
func TestConfigureStorage(t *testing.T) {
	beforeConfigureStorage(t)
	test.Ok(t, ConfigureStorage(bucketName))
}

func beforeTestCreate(t *testing.T) {
	setup(t)
	test.Ok(t, ConfigureStorage(bucketName))
}

func TestCreate(t *testing.T) {
	// test initial bucket create
	beforeTestCreate(t)
	ctx := context.Background()
	test.Ok(t, Create(ctx))

	// test duplicate bucket creation
	test.Ok(t, Create(ctx))
}

func beforeTestDelete(t *testing.T) {
	setup(t)
	test.Ok(t, ConfigureStorage(bucketName))

	ctx := context.Background()
	test.Ok(t, Create(ctx))
}

func TestDelete(t *testing.T) {
	beforeTestDelete(t)
	ctx := context.Background()
	test.Ok(t, Delete(ctx))
}

// func TestUploadFile(t *testing.T) {

// }
