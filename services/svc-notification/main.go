package main

import (
	"fmt"
	"log"
	"net"
	"ops-monorepo/services/svc-notification/config"
	"ops-monorepo/services/svc-notification/internal"
)

func main() {

	// load config
	config, _ := config.LoadConfig(".env")

	// init server
	grpc := internal.NewGrpcServer(config)

	// register implementation
	grpc.Register()

	// create tcp listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// serve
	grpc.Log.Info(fmt.Sprintf("grpc server starting on port: %v", config.Port))
	log.Fatal(grpc.Server.Serve(lis))
}