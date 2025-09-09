package internal

import (
	"net/http"
	"strconv"
	"taskmaster/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// handle GET /api/v1/users/:id - return the user with the matching id
func HandleGetUserByID(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		extractedId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			respondError(c, http.StatusBadRequest, "error extracting id from params")
			logger.API.Println("Error extracting ID from params")
			return
		}

		// get user from supabase
		user, err := s.GetUserById(extractedId)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error getting user")
			logger.DB.Printf("error getting user from db: %v\n", err)
			return
		} else if user == nil {
			// error == user == nil - means no such user exists
			respondError(c, http.StatusNotFound, "no such user exists")
			logger.DB.Printf("No user exists with id %v", extractedId)
			return
		}

		c.JSON(http.StatusFound, user)
	}
}

// handle POST /api/v1/users/ - create new user
func HandlePostUser(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		newUser := User{}
		// bind the json of the request to a User object
		if err := c.BindJSON(&newUser); err != nil {
			respondError(c, http.StatusBadRequest, "error creating task from request params")
			logger.API.Printf("error creating task from params: %v", err)
			return
		}

		// set the CreatedAt date
		newUser.CreatedAt = time.Now()

		_, err := s.CreateNewUser(&newUser)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error making new user")
			logger.DB.Printf("Error making new user: %v\n", err)
			return
		}
		c.JSON(http.StatusCreated, newUser)
	}
}
