package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/nurmuhammaddeveloper/medium_api_gateway/api"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/config"
	grpcPkg "github.com/nurmuhammaddeveloper/medium_api_gateway/pkg/grpc_client"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/pkg/logger"
)

func main() {
	cfg := config.Load(".")

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}

	logger := logger.New()

	apiServer := api.New(&api.RouterOptions{
		Cfg:        &cfg,
		GrpcClient: grpcConn,
		Logger: logger,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
