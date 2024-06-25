package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"template/internal/config"
	"template/internal/server"
)

const cfgPath = "./config/config.yml"

func Start() {
	cfg, err := config.GetConfig(cfgPath)
	if err != nil {
		fmt.Println("config err:", err)
		return
	}
	app := server.NewApp(cfg)
	if err := app.Initialize(); err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		cancel()
	}()
	app.Run(ctx)
}
