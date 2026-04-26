package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lwlee2608/go-reference/internal/api/http/handler"
	"github.com/lwlee2608/go-reference/internal/api/http/middleware"
	"github.com/lwlee2608/go-reference/internal/db/sqlc"
)

type Config struct {
	Port uint
}

type Services struct {
	Queries *sqlc.Queries
}

func SetupRoute(engine *gin.Engine, srvs *Services) {
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.ErrorHandler())

	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(srvs.Queries)

	engine.GET("/health", healthHandler.Check)

	apis := engine.Group("/api/v1")
	{
		users := apis.Group("/users")
		users.POST("", userHandler.Create)
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.Get)
		users.PATCH("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}
}
