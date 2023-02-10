package main

import (
	"log"
	"net"
	"time"

	userpb "github.com/johnnrails/ddd_go/grpc/gen/go/user/v1"
	wearablepb "github.com/johnnrails/ddd_go/grpc/gen/go/wearable/v1"
	"github.com/johnnrails/ddd_go/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userServer := &server.UserService{}
	wearableServer := &server.WearableServer{}

	healthServer := health.NewServer()

	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING

			if time.Now().Second()%2 == 0 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}

			healthServer.SetServingStatus(userpb.UserService_ServiceDesc.ServiceName, status)
			healthServer.SetServingStatus("", status)

			time.Sleep(1 * time.Second)
		}
	}()

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(userpb.UserService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	userpb.RegisterUserServiceServer(grpcServer, userServer)
	wearablepb.RegisterWearableServiceServer(grpcServer, wearableServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	grpcServer.Serve(lis)
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
// type UserServiceServer interface {
//	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
//	mustEmbedUnimplementedUserServiceServer()
// }
