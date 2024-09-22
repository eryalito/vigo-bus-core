package main

import (
	_ "github.com/eryalito/vigo-bus-core/docs" // This is required for the generated docs to be included
	"golang.org/x/time/rate"

	"github.com/eryalito/vigo-bus-core/internal/config"
	"github.com/eryalito/vigo-bus-core/internal/handlers"
	"github.com/eryalito/vigo-bus-core/internal/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title My Microservice API
// @version 1.0
// @description This is a sample microservice API.
// @host localhost:8080
// @BasePath /

func main() {
	config.Init()
	r := gin.Default()

	// Swagger endpoint (no auth middleware)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply rate limiter middleware to all routes
	r.Use(middleware.RateLimiterMiddleware(rate.Limit(config.RateLimiter.Limit), config.RateLimiter.Burst))

	// API endpoints with auth middleware
	api := r.Group("/api")
	// api.Use(middleware.AuthMiddleware)
	// {
	api.GET("/stops", handlers.ListStops)
	api.GET("/stops/:stop_number", handlers.GetStop)
	api.GET("/stops/find", handlers.FindStops)
	api.GET("/stops/find/location", handlers.FindStopsByLocation)
	api.GET("/lines", handlers.ListStops)

	api.GET("/users/:provider/:uuid", handlers.GetUser)
	api.POST("/users/:provider/:uuid", handlers.CreateUser)
	api.POST("/users/:provider/:uuid/favorite_stops/:stop_number", handlers.AddFavoriteStopToIdentity)
	api.DELETE("/users/:provider/:uuid/favorite_stops/:stop_number", handlers.RemoveFavoriteStopFromIdentity)
	// }

	r.Run(":" + config.Port)
}
