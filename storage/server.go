// Package main provides ...
package main

import (
	"fmt"
	"log"
	"net"

	storage "github.com/evanharmon/eph-music-micro/storage"
	"github.com/evanharmon/eph-music-micro/storage/proto/storagepb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	storageSrv := storage.New("glossy-fastness-216519", "test")
	s := grpc.NewServer()
	storagepb.RegisterStorageServer(s, &storageSrv)
}
