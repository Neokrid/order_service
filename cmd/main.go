package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	configs "github.com/Neokrid/order_service/config"
	"github.com/Neokrid/order_service/internal/app"
)

const confDir = "./configs/main.yaml"

func main() {
	cfg, err := configs.NewConfig(confDir, ".env.orders")
	if err != nil {
		log.Fatalf("ошибка конфига: %s", err)
	}
	cnt := app.NewContainer(cfg)
	if err := cnt.Start(context.Background()); err != nil {
		log.Fatalf("ошибка старта: %s", err.Error())
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("%s", "app - Start - signal: "+s.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cnt.Stop(shutdownCtx); err != nil {
		log.Printf("ошибка стопа: %v", err)
	}
}
