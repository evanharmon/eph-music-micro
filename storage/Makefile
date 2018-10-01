
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --go_out=plugins=grpc:. proto/storagepb/storage.proto

build: proto
	go build -o storage storage.go

test:
	go test -v ./... -cover

docker:
	docker build . -t storage:latest
