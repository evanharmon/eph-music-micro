
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/storage/storage.proto

build: proto
	go build -o storage-api storage.go

test:
	go test -v ./... -cover

docker:
	docker build . -t storage-api:latest
