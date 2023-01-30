package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Kyejoon-Lee/grpc-server/config"
	"github.com/Kyejoon-Lee/grpc-server/ent"

	_ "github.com/lib/pq"
)

func connectionURL(cfg *config.Config) (string, error) {
	switch cfg.DBAdapter {
	case "mysql":
		//v := &url.Values{}
		return "", errors.New("Not support mysql")
	case "postgres":
		var host string
		var port string
		var err error
		if strings.Contains(cfg.DBHost, ":") {
			host, port, err = net.SplitHostPort(cfg.DBHost)
			if err != nil {
				return "", err
			}
		} else {
			host = cfg.DBHost
			port = cfg.DBPort
		}
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			host, port, cfg.DBUsername, cfg.DBPassword, cfg.DBName, "disable", cfg.TimeZone,
		), nil
	}
	return "", errors.New("unknown adapter")
}

func AutoMigrate(cfg *config.Config) *ent.Client {

	var (
		url string
		err error
	)

	if url, err = connectionURL(cfg); err != nil {
		panic(err)
	}
	log.Println(url)
	client, err := ent.Open("postgres", url)
	if err != nil {
		log.Fatalf("DB connection failed : %v", err)
	}
	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Schema create failed : %v", err)
	}
	return client
}
