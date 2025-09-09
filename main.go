package main

import (
	"log"
	i "taskmaster/internal"

	gin "github.com/gin-gonic/gin"
)

func main() {
	// Connect to supabase bucket
	sbClient, err := i.CreateSupabaseClient()
	if err != nil {
		log.Fatalf("Error creating supabase client: %v", err)
	}

	// start gin client
	ginClient := gin.Default()
	if err != nil {
		log.Fatalf("Error starting gin client: %v", err)
	}

	tasksGroup := ginClient.Group("/api/v1/tasks")
	{
		tasksGroup.GET("/byuser/:id", func(c *gin.Context) { i.HandleGetTasksByUserId(c)(sbClient) })
		tasksGroup.POST("/", func(c *gin.Context) { i.HandlePostTask(c)(sbClient) })
	}
	ginClient.Run(":3000")
}
