package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kyejoon-Lee/grpc-gateway/api"
	"github.com/Kyejoon-Lee/grpc-gateway/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RestServer struct {
	server *http.Server
}

var (
	cfg = config.GetConfig()
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

	r.GET("/login", api.Login)

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

	log.Infof("start rest server on %s port", cfg.GatewayPort)

}

func (s *RestServer) ShutdownWebServer(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
