package main

import (
	"net"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "grpc-mongodb-crud/proto"
	db "grpc-mongodb-crud/core/config"
	bookstore "grpc-mongodb-crud/core/handlers"
)

var (
	port = ":8000"
)

func startGRPCServer() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookstoreServer(s, &bookstore.Server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	s.Serve(ln)
}

func main(){
	go db.Init()
	startGRPCServer()
}
