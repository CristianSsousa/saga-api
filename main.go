// @title           Saga API
// @version         1.0
// @description     Media tracker backend — movies, TV, games, books, manga
// @termsOfService  http://swagger.io/terms/

// @contact.name   Cristian Sousa
// @contact.email  cristian.sousa365@outlook.com

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	deliveryhttp "github.com/CristianSsousa/saga-api/internal/delivery/http"
	"github.com/CristianSsousa/saga-api/internal/delivery/http/handler"
	"github.com/CristianSsousa/saga-api/internal/infrastructure/apis"
	"github.com/CristianSsousa/saga-api/internal/infrastructure/db"
	"github.com/CristianSsousa/saga-api/internal/repository"
	"github.com/CristianSsousa/saga-api/internal/usecase"
)

func main() {
	_ = godotenv.Load() // loads .env if present; ignored in production

	// Config from environment
	databaseURL := mustEnv("DATABASE_URL")
	jwtSecret := mustEnv("JWT_SECRET")
	tmdbKey := os.Getenv("TMDB_API_KEY")
	rawgKey := os.Getenv("RAWG_API_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Database
	pool, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	mediaCacheRepo := repository.NewMediaCacheRepo(pool)
	libraryRepo := repository.NewLibraryRepo(pool)

	// External search service
	searchSvc := apis.NewAggregatedSearchService(apis.Config{
		TMDBApiKey: tmdbKey,
		RAWGApiKey: rawgKey,
	})

	// Use cases
	authUC := usecase.NewAuthUsecase(userRepo, jwtSecret)
	searchUC := usecase.NewSearchUsecase(searchSvc, mediaCacheRepo)
	libraryUC := usecase.NewLibraryUsecase(libraryRepo, mediaCacheRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC, userRepo)
	searchHandler := handler.NewSearchHandler(searchUC)
	libraryHandler := handler.NewLibraryHandler(libraryUC)

	// Router
	router := deliveryhttp.NewRouter(authHandler, searchHandler, libraryHandler, jwtSecret)
	engine := router.Setup()

	log.Printf("Saga API starting on :%s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
