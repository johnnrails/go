package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/johnnrails/ddd_go/third_ddd_go/common/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedHelloServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	fmt.Println("Running gRPC Server")

	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &Server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
