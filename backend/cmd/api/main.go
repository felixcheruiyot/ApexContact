package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apexcontact/backend/internal/api"
	"github.com/apexcontact/backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Pretty logging in development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	app, cleanup, err := api.NewApp(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialise application")
	}
	defer cleanup()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", ":"+cfg.Port).Msg("ApexContact API starting")
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down...")
	_ = app.ShutdownWithContext(context.Background())
}
