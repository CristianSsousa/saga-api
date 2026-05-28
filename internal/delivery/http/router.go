package http

import (
	"github.com/CristianSsousa/saga-api/internal/delivery/http/handler"
	"github.com/CristianSsousa/saga-api/internal/delivery/http/middleware"
	_ "github.com/CristianSsousa/saga-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine        *gin.Engine
	authHandler   *handler.AuthHandler
	searchHandler *handler.SearchHandler
	libraryHandler *handler.LibraryHandler
	jwtSecret     string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	searchHandler *handler.SearchHandler,
	libraryHandler *handler.LibraryHandler,
	jwtSecret string,
) *Router {
	return &Router{
		engine:        gin.New(),
		authHandler:   authHandler,
		searchHandler: searchHandler,
		libraryHandler: libraryHandler,
		jwtSecret:     jwtSecret,
	}
}

func (r *Router) Setup() *gin.Engine {
	e := r.engine

	// Global middleware
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(corsMiddleware())
	e.Use(middleware.RateLimit())

	// Swagger
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.GET("/me", middleware.Auth(r.jwtSecret), r.authHandler.Me)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.Auth(r.jwtSecret))
	{
		protected.GET("/search", r.searchHandler.Search)

		lib := protected.Group("/library")
		{
			lib.GET("", r.libraryHandler.List)
			lib.POST("", r.libraryHandler.Add)
			lib.PUT("/:id", r.libraryHandler.Update)
			lib.DELETE("/:id", r.libraryHandler.Delete)
		}
	}

	return e
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
