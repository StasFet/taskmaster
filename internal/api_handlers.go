package internal

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// handle /api/v1/tasks/byuser/:id
func HandleGetTasksByUserId(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		// extract user id from request parameters
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Error extracting user id from request params",
			})
			log.Printf("Error extracting user id from request params: %v\n", err)
			return
		}

		// get tasks from supabase
		tasks, err := s.GetTasksByUserId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error retrieving tasks for user",
			})
			log.Printf("Error retruieving tasks for user: %v\n", err)
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
