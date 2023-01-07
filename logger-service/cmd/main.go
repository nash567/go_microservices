package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/loger-service/cmd/app"
)

const (
	defaultConfPath = "./config.yml"
)

func main() {
	var configFiles string
	flag.StringVar(&configFiles, "config", defaultConfPath, "comma separated list of config files to load")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := new(app.Application)
	application.Init(ctx, configFiles)
	application.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// locking till signal received
	<-sigterm
	// start graceful shutdown
	fmt.Println("shutting down")
	application.Stop(ctx)

}
