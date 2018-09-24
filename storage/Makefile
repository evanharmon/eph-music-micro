
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc proto/storagepb/storage.proto --go_out=. --plugin=grpc

build: proto
	go build -o storage-api storage.go

test:
	go test -v ./... -cover

docker:
	docker build . -t storage-api:latest
