GOPATH:=$(shell go env GOPATH)

.PHONY: proto mocks test docker


proto:
	protoc --go_out=plugins=grpc:. proto/storagepb/storage.proto

protomock:
	mockgen \
	-package mocks \
	github.com/evanharmon/eph-music-micro/storage/proto/storagepb \
	StorageClient,Storage_UploadFileClient,Storage_UploadFileServer \
	> core/mocks/mock_storagepb.go

mocks:
	go generate ./...

build: proto protomock mocks
	go build ./...

test:
	go test -v ./... -cover

docker:
	docker build . -t storage:latest
