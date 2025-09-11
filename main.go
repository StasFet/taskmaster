package main

import (
	"os"
	i "taskmaster/internal"
	"taskmaster/logger"
	"time"
	"unicode/utf8"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load("./.env"); err != nil {
		logger.GEN.Fatalf("Error loading .env: %v\n", err)
	}

	// Connect to supabase bucket
	sbClient, err := i.CreateSupabaseClient()
	if err != nil {
		logger.GEN.Fatalf("Error creating supabase client: %v\n", err)
	}

	// start gin client
	ginClient := gin.Default()
	ginClient.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// declare handler groups/routing
	tasksGroup := ginClient.Group("/api/v1/tasks")
	{
		tasksGroup.GET("/", i.JWTValidatorMiddleware(), i.HandleGetTasksByUUID(sbClient))
		tasksGroup.POST("/", i.JWTValidatorMiddleware(), i.HandlePostTask(sbClient))
	}

	usersGroup := ginClient.Group("/api/v1/users")
	{
		usersGroup.GET("/", i.JWTValidatorMiddleware(), i.HandleGetUser(sbClient))
		usersGroup.PUT("/", i.JWTValidatorMiddleware(), i.HandlePutUser(sbClient))
	}

	port := os.Getenv("PORT")
	if (utf8.RuneCountInString(port)) == 0 {
		port = "3000"
	}
	ginClient.Run(":" + port)
}
