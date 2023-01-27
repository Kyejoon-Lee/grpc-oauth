package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Kyejoon-Lee/grpc-gateway/config"
	"github.com/Kyejoon-Lee/grpc-gateway/ent/proto/entpb"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = config.GetConfig()
)

type GrpcCLI struct {
	Connection           *grpc.ClientConn
	UserServiceClientCLI entpb.UserServiceClient
}

func (g *GrpcCLI) StartGrpcConnection() {
	//For making connection to GRPC servers you must make grpc dial
	//to several servers.

	g.Connection, _ = grpc.Dial(cfg.ServerIP+":"+cfg.ServerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	g.UserServiceClientCLI = entpb.NewUserServiceClient(g.Connection)

	log.Infof("start gRPC client to %s:%s server", cfg.ServerIP, cfg.ServerPort)

}

func (g *GrpcCLI) ShutdownGrpcConnection() error {
	return g.Connection.Close()
}
