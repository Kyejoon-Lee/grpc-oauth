package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Kyejoon-Lee/grpc-gateway/config"
	"github.com/Kyejoon-Lee/grpc-gateway/ent/proto/entpb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RestServer struct {
	server *http.Server
}

var (
	cfg = config.GetConfig()
	cli entpb.UserServiceClient
)

func (s *RestServer) StartGatewayServer() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           1,
	}))

	//For making connection to GRPC servers you must make grpc dial
	//to several servers.
	conn, _ := grpc.Dial(cfg.ServerIP+":"+cfg.ServerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	cli = entpb.NewUserServiceClient(conn)

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.GatewayPort),
		Handler:           r,
		ReadHeaderTimeout: 30 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (s *RestServer) ShutdownWebServer(ctx context.Context) error {

	return s.server.Shutdown(ctx)
}
