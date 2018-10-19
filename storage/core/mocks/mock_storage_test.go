package mocks_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/evanharmon/eph-music-micro/storage/core/mocks"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"github.com/golang/mock/gomock"
)

func TestListBuckets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockStorageClient(ctrl)
	req := &pb.ListBucketsRequest{Project: &pb.Project{Id: projectId}}
	res := &pb.ListBucketsResponse{
		Buckets: []*pb.Bucket{
			&pb.Bucket{Name: "evan-terraform-admin"},
			&pb.Bucket{Name: bucketName},
		},
	}
	mock.EXPECT().ListBuckets(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testListBuckets(t, mock)
}

func testListBuckets(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ListBucketsRequest{
		Project: &pb.Project{Id: projectId},
	}
	r, err := client.ListBuckets(ctx, req)
	if err != nil || r.Buckets == nil {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Buckets)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockStorageClient(ctrl)
	req := &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	res := &pb.CreateResponse{Result: "success"}
	mock.EXPECT().Create(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testCreate(t, mock)
}

func testCreate(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	r, err := client.Create(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockStorageClient(ctrl)
	req := &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	res := &pb.DeleteResponse{Result: "success"}
	mock.EXPECT().Delete(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testDelete(t, mock)
}

func testDelete(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	r, err := client.Delete(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockStorageClient(ctrl)
	stream := mocks.NewMockStorage_UploadFileClient(ctrl)
	// Mock All stream Expected Calls On Receiver Functions
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	stream.EXPECT().CloseSend().Return(nil)
	mock.EXPECT().UploadFile(
		gomock.Any(),
	).Return(stream, nil)
	testUploadFile(t, mock)
}

func testUploadFile(t *testing.T, client pb.StorageClient) {
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Error(err)
	}
	fpath, err := filepath.Abs(fileName)
	if err != nil {
		t.Error(err)
	}
	req := &pb.UploadFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		Chunk:   &pb.Chunk{Content: []byte{}},
		File:    &pb.File{Name: fileName, Path: fpath},
	}
	if err := stream.Send(req); err != nil {
		t.Error(err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Error(err)
	}
}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockStorageClient(ctrl)
	req := &pb.DeleteFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		File:    &pb.File{Name: fileName},
	}
	res := &pb.DeleteFileResponse{Result: "success"}
	mock.EXPECT().DeleteFile(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testDeleteFile(t, mock)
}

func testDeleteFile(t *testing.T, client pb.StorageClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DeleteFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		File:    &pb.File{Name: fileName},
	}
	r, err := client.DeleteFile(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}
