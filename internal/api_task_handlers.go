package internal

import (
	"net/http"
	"taskmaster/logger"
	"time"

	db "taskmaster/internal/database"
	model "taskmaster/internal/models"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, code int, text string) {
	c.JSON(code, map[string]any{
		"error_message": text,
	})
}

// handle GET /api/v1/tasks/byuser/:id
func HandleGetTasksByUUID(s *db.SupabaseClient) gin.HandlerFunc {
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
func HandlePostTask(s *db.SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		newTask := model.Task{}
		// bind the json of the request to a Task object
		if err := c.BindJSON(&newTask); err != nil {
			respondError(c, http.StatusBadRequest, "error creating task from request parameters")
			logger.API.Printf("Error binding request params to json: %v\n", err)
			return
		}

		// ensure the CreatedAt date is correct
		newTask.CreatedAt = time.Now()
		newTask.OwnerUUID = c.GetString("validated_uuid")
		newTask.Status = "INCOMPLETE"

		returnedTask, err := s.CreateNewTask(&newTask)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error creating new task on supabase")
			logger.DB.Printf("Error creating new task into supabase: %v\n", err)
			return
		}
		c.JSON(http.StatusCreated, *returnedTask)
	}
}

// handlePutTask
func HandlePutTask(s *db.SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		updatedTask := model.Task{}
		if err := c.BindJSON(&updatedTask); err != nil {
			respondError(c, http.StatusBadRequest, "error binding updated task json")
			logger.API.Printf("Error binding json for PUT task request: %v\n", err)
			return
		}

		updatedTask.OwnerUUID = c.GetString("validated_uuid")

		newTask, err := s.UpdateTask(&updatedTask)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error updating task")
			logger.DB.Printf("Error updating task: %v\n", err)
			return
		}
		c.JSON(http.StatusOK, newTask)
	}
}

// handle DELETE task
func HandleDeleteTask(s *db.SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		uuid := c.GetString("validated_uuid")
		err := s.DeleteTask(id, uuid)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error deleting task")
			logger.DB.Printf("Error deleting task with id %v: %v", id, err)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}
