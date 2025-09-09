package internal

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, code int, text string) {
	c.JSON(code, map[string]any{
		"error_message": text,
	})
}

// handle GET /api/v1/tasks/byuser/:id
func HandleGetTasksByUserId(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		// extract user id from request parameters
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			respondError(c, http.StatusBadRequest, "error extracting user id from request params")
			log.Printf("Error extracting user id from request params: %v\n", err)
			return
		}

		// get tasks from supabase
		tasks, err := s.GetTasksByUserId(userId)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "Error retrieving tasks for user")
			log.Printf("Error retruieving tasks for user: %v\n", err)
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}

// handle POST /api/v1/tasks - Create new task
func HandlePostTask(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		newTask := Task{}
		if err := c.BindJSON(&newTask); err != nil {
			respondError(c, http.StatusBadRequest, "error creating task from request parameters")
			log.Printf("error binding request params to json: %v\n", err)
			return
		}

		// ensure the createdat date is correct
		newTask.CreatedAt = time.Now()

		_, err := s.CreateNewTask(&newTask)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error creating new task on supabase")
			log.Printf("error creating new task into supabase: %v\n", err)
			return
		}
	}
}
