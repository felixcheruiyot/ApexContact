package api

import (
	"context"
	"fmt"

	"github.com/livestreamify/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func setupInfrastructure(cfg *config.Config) (*pgxpool.Pool, *redis.Client, func(), error) {
	// PostgreSQL
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("postgres: %w", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		return nil, nil, nil, fmt.Errorf("postgres ping: %w", err)
	}
	log.Info().Msg("postgres connected")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("redis ping: %w", err)
	}
	log.Info().Msg("redis connected")

	cleanup := func() {
		db.Close()
		_ = rdb.Close()
		log.Info().Msg("infrastructure connections closed")
	}

	return db, rdb, cleanup, nil
}
