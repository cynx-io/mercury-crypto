package main

import (
	"context"
	"log"
	"mercury/internal/app"
	"mercury/internal/pkg/logger"
)

func main() {
	log.Println("Starting mercury")
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	log.Println("Initializing App")
	application, err := app.NewApp("config")
	if err != nil {
		panic(err)
	}

	logger.Info("Creating servers")
	servers, err := application.NewServers()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting servers")
	if err := servers.Start(); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
