package server

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/johnnrails/ddd_go/third_ddd_go/common/logs"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	logger := logrus.New()
	logs.SetFormatter(logger)
	logger.SetLevel(logrus.WarnLevel)
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

func newGrpcServer() *grpc.Server {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	return grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
			),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)
}

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	grpcServer := newGrpcServer()
	registerServer(grpcServer)

	listen, err := net.Listen("tpc", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.WithField("grpcEndpoint", ":8080").Info("Starting: gRPC Listener")
	logrus.Fatal(grpcServer.Serve(listen))
}
