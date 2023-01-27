package app

import (
	log "github.com/sirupsen/logrus"

	"net"

	"google.golang.org/grpc"

	"github.com/Kyejoon-Lee/grpc-server/config"
	"github.com/Kyejoon-Lee/grpc-server/ent/proto/entpb"
)

type GrpcServer struct {
	Server *grpc.Server
}

type UserServer struct {
	entpb.UserServiceServer
}

var cfg = config.GetConfig()

func (s *GrpcServer) StartGrpcServer() {
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	entpb.RegisterUserServiceServer(s.Server, &UserServer{})

	log.Printf("start gRPC server on %s port", cfg.ServerPort)
	go func() {
		if err := s.Server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
}

func (s *GrpcServer) ShutdownGrpcServer() {
	s.Server.GracefulStop()
}
