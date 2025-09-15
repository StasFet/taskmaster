package main

import (
	"os"
	i "taskmaster/internal"
	db "taskmaster/internal/database"
	sec "taskmaster/internal/security"
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
	sbClient, err := db.CreateSupabaseClient()
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
		MaxAge:           6 * time.Hour,
	}), sec.JWTValidatorMiddleware())

	// declare handler groups/routing
	tasksGroup := ginClient.Group("/api/v1/tasks")
	{
		tasksGroup.GET("/", i.HandleGetTasksByUUID(sbClient))
		tasksGroup.POST("/", i.HandlePostTask(sbClient))
		tasksGroup.PUT("/", i.HandlePutTask(sbClient))
		tasksGroup.DELETE("/:id", i.HandleDeleteTask(sbClient))
	}

	usersGroup := ginClient.Group("/api/v1/users")
	{
		usersGroup.GET("/", i.HandleGetUser(sbClient))
		usersGroup.PUT("/", i.HandlePutUser(sbClient))
	}

	port := os.Getenv("PORT")
	if (utf8.RuneCountInString(port)) == 0 {
		port = "3000"
	}
	ginClient.Run(":" + port)
}
