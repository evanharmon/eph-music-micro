// +build integration

package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	helper "github.com/evanharmon/eph-music-micro/helper"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storage"
	"github.com/google/uuid"
)

func testProjectID(t *testing.T) string {
	t.Helper()
	id, err := helper.GetEnv("GOOGLE_PROJECT_ID")
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("ProjectID must not be an empty string")
	}
	return id
}

func testBucketName(t *testing.T, id string) string {
	t.Helper()
	return fmt.Sprintf("test-eph-music-%s", id)
}

func TestConfigure(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if _, err := configure(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestNew(t *testing.T) {
	id := testProjectID(t)
	bid := uuid.New().String()
	name := testBucketName(t, bid)

	t.Run("success", func(t *testing.T) {
		if _, err := New(id, name); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("empty project id", func(t *testing.T) {
		client, err := New("", name)
		if client != nil || err == nil {
			t.Error("Empty string project ID should return an error")
		}
	})
	t.Run("empty bucket name", func(t *testing.T) {
		client, err := New(id, "")
		if client != nil || err == nil {
			t.Error("Empty string bucket name should return an error")
		}
	})
}

func testNew(t *testing.T, id string, name string) Service {
	t.Helper()
	svc, err := New(id, name)
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestCreate(t *testing.T) {
	id := testProjectID(t)
	bid := uuid.New().String()
	name := testBucketName(t, bid)

	tests := map[string]struct {
		name string
	}{
		"create success":           {name},
		"duplicate bucket success": {name},
	}

	t.Parallel()
	t.Run("success", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// SETUP
				s := testNew(t, id, test.name)
				// TEST
				p := &pb.Project{Id: id}
				if err := s.Create(context.Background(), p); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		// SETUP
		s := testNew(t, id, "test") // always taken
		// TEST
		p := &pb.Project{Id: id}
		if err := s.Create(context.Background(), p); err == nil {
			t.Fatal("should throw error on duplicate bucket")
		}
	})

	testDelete(t, id, name)
}

func testCreate(t *testing.T, id string, name string) {
	t.Helper()
	s := testNew(t, id, name)
	p := &pb.Project{Id: id}
	if err := s.Create(context.Background(), p); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	id := testProjectID(t)
	bid := uuid.New().String()
	name := testBucketName(t, bid)
	s := testNew(t, id, name)
	testCreate(t, id, name)

	t.Parallel()
	t.Run("success", func(t *testing.T) {
		if err := s.Delete(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func testDelete(t *testing.T, id string, name string) {
	t.Helper()
	testCreate(t, id, name)
	s := testNew(t, id, name)
	if err := s.Delete(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile(t *testing.T) {
	id := testProjectID(t)
	fid := uuid.New().String()
	name := testBucketName(t, fid)
	testCreate(t, id, name)
	s := testNew(t, id, name)
	path, err := filepath.Abs("./testdata/upload-file.txt")
	if err != nil {
		t.Error(err)
	}

	t.Parallel()
	t.Run("success", func(t *testing.T) {
		if err := s.UploadFile(context.Background(), path); err != nil {
			t.Fatal(err)
		}
	})

	testDeleteFile(t, id, name, "upload-file.txt")
	testDelete(t, id, name)
}

func testUploadFile(t *testing.T, id string, name string, fname string) {
	t.Helper()
	s := testNew(t, id, name)
	testCreate(t, id, name)
	path, err := filepath.Abs(fname)
	if err != nil {
		t.Error(err)
	}
	if err := s.UploadFile(context.Background(), path); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFile(t *testing.T) {
	id := testProjectID(t)
	fid := uuid.New().String()
	name := testBucketName(t, fid)
	fname := "./testdata/upload-file.txt"
	testUploadFile(t, id, name, fname)
	s := testNew(t, id, name)

	t.Parallel()
	t.Run("success", func(t *testing.T) {
		err := s.DeleteFile(context.Background(), "upload-file.txt")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func testDeleteFile(t *testing.T, id string, name string, fname string) {
	t.Helper()
	s := testNew(t, id, name)
	if err := s.DeleteFile(context.Background(), fname); err != nil {
		t.Fatal(err)
	}
}
