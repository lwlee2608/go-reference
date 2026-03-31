package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL      string `mask:"true"`
	Schema   string
	MaxConns int32
	MinConns int32
}

func InitDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	schema := cfg.Schema
	if schema != "" {
		if poolConfig.ConnConfig.RuntimeParams == nil {
			poolConfig.ConnConfig.RuntimeParams = map[string]string{}
		}
		poolConfig.ConnConfig.RuntimeParams["search_path"] = schema
		slog.Info("Setting search_path for connection pool", "schema", schema)

		poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			_, err := conn.Exec(ctx, fmt.Sprintf("SET search_path TO %s", pgx.Identifier{schema}.Sanitize()))
			if err != nil {
				slog.Warn("Failed to set search_path in AfterConnect", "error", err)
				return err
			}
			slog.Debug("Set search_path for new connection", "schema", schema)
			return nil
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Connected to PostgreSQL")

	return pool, nil
}
