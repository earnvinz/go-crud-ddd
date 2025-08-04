package main

import (
	"log"
	"os"
	database "test-go/db"
	customer "test-go/internal/customer"
	healthcheck "test-go/internal/health-check"

	_ "test-go/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title Backend-Go-API
// @version 1.0
// @BasePath /api/v1
func main() {

	err := godotenv.Load()
	db, err := database.ConnectPostgres()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	router := setupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Failed to run server:", err)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		customer.RegisterRoutes(apiV1, db)
		healthcheck.RegisterRoutes(apiV1, db)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
