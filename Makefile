GOPATH:=$(shell go env GOPATH)

build:
	go build -o eph-music-micro main.go
