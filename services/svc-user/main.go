package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"ops-monorepo/services/svc-user/config"
	"ops-monorepo/services/svc-user/internal"
)

func main() {

	// load config
	config, _ := config.LoadConfig(".env")

	// init servers
	grpcServer := internal.NewGrpcServer(config)
	httpServer := internal.NewHTTPServer(config)

	// register implementations
	grpcServer.Register()
	httpServer.Register()

	// start HTTP server in a goroutine
	go func() {
		httpServer.Log.Info(fmt.Sprintf("HTTP server starting on port: %v", config.HTTPPort))
		if err := http.ListenAndServe(fmt.Sprintf(":%v", config.HTTPPort), httpServer.Router); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// create tcp listener for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// serve gRPC
	grpcServer.Log.Info(fmt.Sprintf("gRPC server starting on port: %v", config.Port))
	log.Fatal(grpcServer.Server.Serve(lis))
}