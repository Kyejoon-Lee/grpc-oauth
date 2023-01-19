package module

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Kyejoon-Lee/grpc-gateway/config"
)

var serverContextOnce sync.Once
var serverContext context.Context
var serverCancelFunc context.CancelFunc

func ServerContext() (context.Context, func()) {
	serverContextOnce.Do(func() {
		serverContext, serverCancelFunc = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	})
	return serverContext, serverCancelFunc
}

func Config() *config.Config {
	return config.GetConfig()
}
