package internal

import (
	"net/http"
	"taskmaster/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, code int, text string) {
	c.JSON(code, map[string]any{
		"error_message": text,
	})
}

// handle GET /api/v1/tasks/byuser/:id
func HandleGetTasksByUUID(s *SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract user id from request parameters
		uuid := c.GetString("validated_uuid")

		// get tasks from supabase
		tasks, err := s.GetTasksByUUID(uuid)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "Error retrieving tasks for user")
			logger.DB.Printf("Error retrieving tasks for user: %v\n", err)
			return
		}

		c.JSON(http.StatusOK, map[string]any{
			"tasks": tasks,
			"count": len(*tasks),
		})
	}
}

// handle POST /api/v1/tasks - Create new task
func HandlePostTask(s *SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		newTask := Task{}
		// bind the json of the request to a Task object
		if err := c.BindJSON(&newTask); err != nil {
			respondError(c, http.StatusBadRequest, "error creating task from request parameters")
			logger.API.Printf("error binding request params to json: %v\n", err)
			return
		}

		// ensure the CreatedAt date is correct
		newTask.CreatedAt = time.Now()

		_, err := s.CreateNewTask(&newTask)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error creating new task on supabase")
			logger.DB.Printf("error creating new task into supabase: %v\n", err)
			return
		}
		c.JSON(http.StatusCreated, newTask)
	}
}
