package app

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Kyejoon-Lee/grpc-server/config"
	"github.com/Kyejoon-Lee/grpc-server/ent/proto/entpb"
)

type UserServer struct {
	entpb.UserServiceServer
}

var cfg = config.GetConfig()

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	entpb.RegisterUserServiceServer(grpcServer, &UserServer{})

	log.Printf("start gRPC server on %s port", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
