package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	helper "github.com/evanharmon/eph-music-micro/helper"
	testhelper "github.com/evanharmon/eph-music-micro/helper/testhelper"
)

//
var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	bucketName      = ""
)

// Setup sets the necessary vars from ENV exports
func setup(t *testing.T) {
	id, err := helper.GetEnv(envVarProjectID)
	testhelper.Ok(t, err)
	testhelper.Assert(t, len(id) != 0, "Export necessary ENV vars before running tests")
	ProjectID = id
	bucketName = fmt.Sprintf("%s-test", id)
}

// TestConfigureStoraage tests creating a client for the Google Cloud Storage
// API
func TestConfigureStorage(t *testing.T) {
	setup(t)
	testhelper.Ok(t, ConfigureStorage(bucketName))
}

func TestCreate(t *testing.T) {
	setup(t)
	tests := map[string]struct {
		name string
	}{
		"create success":           {bucketName},
		"duplicate bucket success": {bucketName},
		"duplicate bucket failure": {"test"}, // this bucket is always taken in global namespaces
	}

	for k, test := range tests {
		setup(t)
		testhelper.Ok(t, ConfigureStorage(test.name))
		ctx := context.Background()
		if k == "duplicate bucket failure" {
			testhelper.Throws(t, Create(ctx))
		} else {
			bucketName = fmt.Sprintf("%s-test", ProjectID)
			testhelper.Ok(t, Create(ctx))
		}
	}
}

func beforeTestDelete(t *testing.T) {
	setup(t)
	testhelper.Ok(t, ConfigureStorage(bucketName))
	ctx := context.Background()
	testhelper.Ok(t, Create(ctx))
}

func TestDelete(t *testing.T) {
	beforeTestDelete(t)
	ctx := context.Background()
	testhelper.Ok(t, Delete(ctx))
}

func beforeTestUploadFile(t *testing.T) {
	setup(t)
	testhelper.Ok(t, ConfigureStorage(bucketName))
	ctx := context.Background()
	testhelper.Ok(t, Create(ctx))
}

func TestUploadFile(t *testing.T) {
	beforeTestUploadFile(t)
	fpath, err := filepath.Abs("./testdata/upload-file.txt")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	testhelper.Ok(t, UploadFile(ctx, fpath))
}

func beforeTestDeleteFile(t *testing.T) {
	setup(t)
	testhelper.Ok(t, ConfigureStorage(bucketName))
	ctx := context.Background()
	testhelper.Ok(t, Create(ctx))
	TestUploadFile(t)
}

func TestDeleteFile(t *testing.T) {
	beforeTestDeleteFile(t)
	fname := "upload-file.txt"
	ctx := context.Background()
	testhelper.Ok(t, DeleteFile(ctx, fname))
}
