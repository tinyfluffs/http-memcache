package main

import (
	"context"
	. "github.com/tinyfluffs/http/memcache/internal"
	"github.com/tinyfluffs/http/memcache/internal/cache"
	"github.com/caarlos0/env/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	mc := cache.New(cfg.Expiration, cfg.GCInterval, cfg.ChunkCount)
	httpServer := &http.Server{
		Addr: cfg.Address,
	}

	go func() {
		<-signals
		mc.Stop()
		_ = httpServer.Shutdown(context.Background())
	}()

	Run(httpServer, mc)
}
