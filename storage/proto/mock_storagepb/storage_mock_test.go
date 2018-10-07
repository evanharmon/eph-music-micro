package mock_storagepb_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	storagemock "github.com/evanharmon/eph-music-micro/storage/proto/mock_storagepb"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

const projectId = "evan-terraform-admin"
const fname = "upload-file.txt"

type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}

	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestListBuckets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageClient := storagemock.NewMockStorageClient(ctrl)
	req := &pb.ListBucketsRequest{Project: &pb.Project{Id: projectId}}
	mockStorageClient.EXPECT().ListBuckets(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.ListBucketsResponse{
		Buckets: []*pb.Bucket{
			&pb.Bucket{Name: "evan-terraform-admin"},
			&pb.Bucket{Name: "eph-test-music"},
		},
	}, nil)
	testListBuckets(t, mockStorageClient)
}

func testListBuckets(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.ListBuckets(ctx, &pb.ListBucketsRequest{
		Project: &pb.Project{Id: projectId},
	})
	if err != nil || r.Buckets == nil {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Buckets)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageClient := storagemock.NewMockStorageClient(ctrl)
	req := &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
	}
	mockStorageClient.EXPECT().Create(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.CreateResponse{Result: "success"}, nil)
	testCreate(t, mockStorageClient)
}

func testCreate(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Create(ctx, &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
	})
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageClient := storagemock.NewMockStorageClient(ctrl)
	req := &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
	}
	mockStorageClient.EXPECT().Delete(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.DeleteResponse{Result: "success"}, nil)
	testDelete(t, mockStorageClient)
}

func testDelete(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Delete(ctx, &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
	})
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageClient := storagemock.NewMockStorageClient(ctrl)
	// fpath, err := filepath.Abs(fname)
	// if err != nil {
	// t.Error(err)
	// }
	// req := &pb.UploadFileRequest{
	// Project: &pb.Project{Name: projectId},
	// Bucket:  &pb.Bucket{Name: "eph-test-music"},
	// Chunk:   &pb.Chunk{Content: []byte{}},
	// File:    &pb.File{Name: fname, Path: fpath},
	// }
	// mockStorageClient.EXPECT().UploadFile(
	// gomock.Any(),
	// &rpcMsg{msg: req},
	// ).Return(&pb.UploadFileResponse{
	// Message: "success",
	// Code:    pb.UploadStatusCode_Ok,
	// }, nil)
	stream := storagemock.NewMockStorage_UploadFileClient(ctrl)
	// Mock All stream Expected Calls On Receiver Functions
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	stream.EXPECT().CloseSend().Return(nil)
	mockStorageClient.EXPECT().UploadFile(
		gomock.Any(),
	).Return(stream, nil)
	if err := testUploadFile(mockStorageClient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testUploadFile(client pb.StorageClient) error {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		return err
	}
	fpath, err := filepath.Abs(fname)
	if err != nil {
		return err
	}
	req := &pb.UploadFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
		Chunk:   &pb.Chunk{Content: []byte{}},
		File:    &pb.File{Name: fname, Path: fpath},
	}
	if err := stream.Send(req); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}

	return nil
}

func DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageClient := storagemock.NewMockStorageClient(ctrl)
	req := &pb.DeleteFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
		File:    &pb.File{Name: fname},
	}
	mockStorageClient.EXPECT().DeleteFile(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.DeleteFileResponse{Result: "success"}, nil)
	testDeleteFile(t, mockStorageClient)
}

func testDeleteFile(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.DeleteFile(ctx, &pb.DeleteFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: "eph-test-music"},
		File:    &pb.File{Name: fname},
	})
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}
