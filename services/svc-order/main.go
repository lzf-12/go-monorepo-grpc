package main

import (
	"fmt"
	"log"
	"ops-monorepo/services/svc-order/config"
	"ops-monorepo/services/svc-order/internal"
)

func main() {

	// load config
	config, _ := config.LoadConfig(".env")

	// init server
	server := internal.NewHTTPServer(config)

	// register routes
	server.RegisterRoutes()

	// start server
	addr := fmt.Sprintf(":%s", config.Port)
	log.Fatal(server.Start(addr))
}
