package main

import (
	"context"
	"fmt"
	"log"
	"time"

	userpb "github.com/johnnrails/ddd_go/grpc/gen/go/user/v1"
	wearablepb "github.com/johnnrails/ddd_go/grpc/gen/go/wearable/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main_user() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	res, err := client.GetUser(context.Background(), &userpb.GetUserRequest{
		Uuid: "hello",
	})
	if err != nil {
		log.Fatalf("fail to GetUser: %v", err)
	}
	fmt.Printf("%+v\n", res)
}

func main_client_streaming() {

}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := wearablepb.NewWearableServiceClient(conn)
	stream, err := client.ConsumeBeatsPerMinute(context.Background())

	for i := 0; i < 10; i++ {
		fmt.Println("SENDING... ", i)
		err := stream.Send(&wearablepb.ConsumeBeatsPerMinuteRequest{
			Uuid:   "johnn",
			Value:  uint32(i),
			Minute: uint32(i * 2),
		})
		if err != nil {
			log.Fatalln("Send", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Close", err)
	}
	main_client_streaming()
	fmt.Println(resp.GetTotal())
}
