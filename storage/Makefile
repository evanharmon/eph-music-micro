
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --go_out=plugins=grpc:. proto/storagepb/storage.proto

protomock:
		mockgen \
		github.com/evanharmon/eph-music-micro/storage/proto/storagepb \
		StorageClient,Storage_UploadFileClient \
		> proto/mock_storagepb/storage_mock.go

build: proto protomock
	go build -o storage storage.go

test:
	go test -v ./... -cover

docker:
	docker build . -t storage:latest
