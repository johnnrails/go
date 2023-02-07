package main

import (
	"context"
	"log"
	"time"

	"github.com/johnnrails/ddd_go/third_ddd_go/common/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewHelloClient(conn)
	for i := 1; i < 100; i++ {
		start := time.Now()
		res, err := client.SayHello(context.Background(), &pb.HelloRequest{
			Name: "John",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res)
		log.Print(time.Since(start).Milliseconds())
	}

}
