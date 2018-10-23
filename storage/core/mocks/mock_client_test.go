package mocks_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/evanharmon/eph-music-micro/storage/core/mocks"
	pb "github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
)

const (
	bucketName = "test-eph-music"
	fileName   = "upload-file.txt"
	projectId  = "evan-terraform-admin"
)

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

func TestListBucketsClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockClientService(ctrl)
	req := &pb.ListBucketsRequest{Project: &pb.Project{Id: projectId}}
	mock.EXPECT().ListBuckets(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.ListBucketsResponse{
		Buckets: []*pb.Bucket{
			&pb.Bucket{Name: "evan-terraform-admin"},
			&pb.Bucket{Name: bucketName},
		},
	}, nil)
	testListBucketsClient(t, mock)
}

func testListBucketsClient(t *testing.T, c *mocks.MockClientService) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ListBucketsRequest{
		Project: &pb.Project{Id: projectId},
	}
	r, err := c.ListBuckets(ctx, req)
	if err != nil || r.Buckets == nil {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Buckets)
}

func TestCreateClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockClientService(ctrl)
	req := &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	res := &pb.CreateResponse{Result: "success"}
	mock.EXPECT().Create(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testCreateClient(t, mock)
}

func testCreateClient(t *testing.T, c *mocks.MockClientService) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.CreateRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	r, err := c.Create(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestDeleteClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockClientService(ctrl)
	req := &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	res := &pb.DeleteResponse{Result: "success"}
	mock.EXPECT().Delete(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testDeleteClient(t, mock)
}

func testDeleteClient(t *testing.T, c *mocks.MockClientService) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DeleteRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
	}
	r, err := c.Delete(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}

func TestUploadFileClient(t *testing.T) {
	fpath, err := filepath.Abs(fileName)
	if err != nil {
		t.Error(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockClientService(ctrl)
	req := &pb.UploadFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		Chunk:   &pb.Chunk{Content: []byte{}},
		File:    &pb.File{Name: fileName, Path: fpath},
	}
	res := &pb.UploadFileResponse{
		Message: "success",
		Code:    pb.UploadStatusCode_Ok,
	}
	mock.EXPECT().UploadFile(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(res, nil)
	testUploadFileClient(t, mock)
}

func testUploadFileClient(t *testing.T, c *mocks.MockClientService) {
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
	res, err := c.UploadFile(context.Background(), req)
	if err != nil || res.Message == "" {
		t.Error(err)
	}
}

func TestDeleteFileClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockClientService(ctrl)
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
	testDeleteFileClient(t, mock)
}

func testDeleteFileClient(t *testing.T, c *mocks.MockClientService) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DeleteFileRequest{
		Project: &pb.Project{Name: projectId},
		Bucket:  &pb.Bucket{Name: bucketName},
		File:    &pb.File{Name: fileName},
	}
	r, err := c.DeleteFile(ctx, req)
	if err != nil || r.Result == "" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply: ", r.Result)
}
