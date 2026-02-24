package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	internalhttp "github.com/lwlee2608/go-reference/internal/api/http"
	"github.com/lwlee2608/go-reference/internal/db"
	"github.com/lwlee2608/go-reference/internal/db/sqlc"
)

var AppVersion = "dev"

func main() {
	InitConfig()

	slog.Info("go-reference", "version", AppVersion)

	if config.DB.URL == "" {
		panic("db.url is required")
	}

	if err := db.RunMigrations(config.DB.URL, config.DB.Schema); err != nil {
		panic(err)
	}

	dbPool, err := db.InitDB(context.Background(), config.DB.URL, config.DB.Schema)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	services := &internalhttp.Services{
		Queries: sqlc.New(dbPool),
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.Use(gin.Recovery())
	internalhttp.SetupRoute(engine, services)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Http.Port),
		Handler: engine,
	}

	slog.Info("Starting HTTP server", "address", server.Addr)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server error", "error", err)
	}
}
