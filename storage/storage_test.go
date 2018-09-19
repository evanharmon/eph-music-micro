package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	helper "github.com/evanharmon/eph-music-micro/helper"
	"github.com/pborman/uuid"
)

func testProjectID(t *testing.T) string {
	t.Helper()
	projectID, err := helper.GetEnv("GOOGLE_PROJECT_ID")
	if err != nil {
		t.Fatal(err)
	}
	if len(projectID) == 0 {
		t.Fatal("ProjectID must not be an empty string")
	}
	return projectID
}

func testBucketName(t *testing.T, id string) string {
	t.Helper()
	return fmt.Sprintf("test-eph-music-%s", id)
}

func TestConfigure(t *testing.T) {
	uid := uuid.New()
	name := testBucketName(t, uid)

	t.Run("success", func(t *testing.T) {
		if _, err := configure(name); err != nil {
			t.Fatal(err)
		}
	})
}

func TestNew(t *testing.T) {
	id := testProjectID(t)
	uid := uuid.New()
	name := testBucketName(t, uid)

	t.Run("success", func(t *testing.T) {
		if _, err := New(id, name); err != nil {
			t.Fatal(err)
		}
	})
}

func testNew(t *testing.T, id string, name string) Service {
	t.Helper()
	bucketsvc, err := New(id, name)
	if err != nil {
		t.Fatal(err)
	}
	return bucketsvc
}

func TestCreate(t *testing.T) {
	id := testProjectID(t)
	uid := uuid.New()
	name := testBucketName(t, uid)

	tests := map[string]struct {
		name string
	}{
		"create success":           {name},
		"duplicate bucket success": {name},
	}

	t.Run("success", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// SETUP
				bucket := testNew(t, id, test.name)
				// TEST
				ctx := context.Background()
				if err := bucket.Create(ctx); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		// SETUP
		bucket := testNew(t, id, "test") // always taken
		// TEST
		ctx := context.Background()
		if err := bucket.Create(ctx); err == nil {
			t.Fatal("should throw error on duplicate bucket")
		}
	})

	testDelete(t, id, name)
}

func testCreate(t *testing.T, id string, name string) {
	t.Helper()
	bucket := testNew(t, id, name)
	ctx := context.Background()
	if err := bucket.Create(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	id := testProjectID(t)
	uid := uuid.New()
	name := testBucketName(t, uid)
	bucket := testNew(t, id, name)
	testCreate(t, id, name)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		if err := bucket.Delete(ctx); err != nil {
			t.Fatal(err)
		}
	})
}

func testDelete(t *testing.T, id string, name string) {
	t.Helper()
	testCreate(t, id, name)
	bucket := testNew(t, id, name)
	ctx := context.Background()
	if err := bucket.Delete(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile(t *testing.T) {
	id := testProjectID(t)
	uid := uuid.New()
	name := testBucketName(t, uid)
	testCreate(t, id, name)
	bucket := testNew(t, id, name)
	fpath, err := filepath.Abs("./testdata/upload-file.txt")
	if err != nil {
		t.Error(err)
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		if err := bucket.UploadFile(ctx, fpath); err != nil {
			t.Fatal(err)
		}
	})

	testDeleteFile(t, id, name, "upload-file.txt")
	testDelete(t, id, name)
}

func testUploadFile(t *testing.T, id string, name string, fname string) {
	t.Helper()
	bucket := testNew(t, id, name)
	ctx := context.Background()
	testCreate(t, id, name)
	fpath, err := filepath.Abs(fname)
	if err != nil {
		t.Error(err)
	}
	if err := bucket.UploadFile(ctx, fpath); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFile(t *testing.T) {
	id := testProjectID(t)
	uid := uuid.New()
	name := testBucketName(t, uid)
	fname := "./testdata/upload-file.txt"
	testUploadFile(t, id, name, fname)
	bucket := testNew(t, id, name)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		if err := bucket.DeleteFile(ctx, "upload-file.txt"); err != nil {
			t.Fatal(err)
		}
	})
}

func testDeleteFile(t *testing.T, id string, name string, fname string) {
	t.Helper()
	bucket := testNew(t, id, name)
	ctx := context.Background()
	if err := bucket.DeleteFile(ctx, fname); err != nil {
		t.Fatal(err)
	}
}
