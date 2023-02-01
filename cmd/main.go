package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Petr0752/L0/packages/cache"
	"github.com/Petr0752/L0/packages/config"
	"github.com/Petr0752/L0/packages/db"
	"github.com/Petr0752/L0/packages/handler"
	"github.com/Petr0752/L0/packages/nats"
	"github.com/Petr0752/L0/packages/server"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	// load config
	err := config.Load("configs/")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get()
	log.Printf("config: %+v\n", cfg)

	// connect to db
	database, err := db.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// create cache
	cacheImpl, err := cache.New(ctx, database)
	if err != nil {
		log.Fatal(err)
	}

	// connect to nats-streaming
	natsHandler, err := nats.New(cacheImpl)
	if err != nil {
		log.Fatal(err)
	}
	defer natsHandler.Close()

	// start processing messages from nats-streaming
	err = natsHandler.Start()
	if err != nil {
		log.Fatal(err)
	}

	// init HTTP handlers
	h := handler.New(database, cacheImpl)

	// start HTTP server
	srv := server.New(h)
	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	log.Printf("service started %s", cfg.ServerAddress)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Printf("service shutting down")
	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shut down: %v", err)
	}
}
