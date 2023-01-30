/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/Kyejoon-Lee/grpc-server/config"
	"github.com/Kyejoon-Lee/grpc-server/db"
	"github.com/Kyejoon-Lee/grpc-server/internal/app"
	"github.com/Kyejoon-Lee/grpc-server/internal/app/module"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		runServer()
	},
}

func init() {
	config.InitConfig(serverCmd.PersistentFlags())
	err := viper.BindPFlags(serverCmd.PersistentFlags())
	if err != nil {
		panic(fmt.Sprintf("viper - bind flags fail - %v", err))
	}

	rootCmd.AddCommand(serverCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runServer() {
	cfg := config.GetConfig()
	ctx, stop := module.ServerContext()
	defer stop()

	client := db.AutoMigrate(cfg)
	grpcServer := app.GrpcServer{Server: grpc.NewServer()}
	grpcServer.StartGrpcServer(client)
	defer client.Close()
	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.ShutdownGrpcServer()

	log.Info("gRPC server is shutdown")
}
