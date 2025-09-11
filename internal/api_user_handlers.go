package internal

import (
	"net/http"
	"taskmaster/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// handle GET /api/v1/users/:id - return the user with the matching id
func HandleGetUser(s *SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.GetString("validated_uuid")
		name := c.GetString("validated_name")
		email := c.GetString("validated_email")

		// get user from supabase
		user, err := s.GetUserByUUID(uuid)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error getting user")
			logger.DB.Printf("error getting user from db: %v\n", err)
			return
		}

		// no user exists with given uuid
		if user == nil {
			user, err = s.CreateNewUser(&User{Name: name, UUID: uuid, Email: email})
			if err != nil {
				respondError(c, http.StatusInsufficientStorage, "Error creating new user")
				logger.DB.Printf("error creating user for new UUID: %v\n", err)
				return
			}
		}

		c.JSON(http.StatusFound, user)
	}
}

// handle PUT /api/v1/users/ - update user details
func HandlePutUser(s *SupabaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := User{}
		// bind the json of the request to a User object
		if err := c.BindJSON(&newUser); err != nil {
			respondError(c, http.StatusBadRequest, "error creating user from request params")
			logger.API.Printf("error creating user from request params: %v\n", err)
			c.Abort()
			return
		}

		returnedUser, err := s.UpdateUser(c.GetString("validated_uuid"), &newUser)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error updating user details")
			logger.DB.Printf("error updating user data: %v\n", err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, returnedUser)
	}
}

// handle POST /api/v1/users/ - create new user (this is never gonna be used because new users are made upon a GET request with a new UUID)
func HandlePostUser(c *gin.Context) func(s *SupabaseClient) {
	return func(s *SupabaseClient) {
		newUser := User{}
		// bind the json of the request to a User object
		if err := c.BindJSON(&newUser); err != nil {
			respondError(c, http.StatusBadRequest, "error creating user from request params")
			logger.API.Printf("error creating user from params: %v", err)
			c.Abort()
			return
		}

		// set the CreatedAt date
		newUser.CreatedAt = time.Now()

		_, err := s.CreateNewUser(&newUser)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "error making new user")
			logger.DB.Printf("Error making new user: %v\n", err)
			c.Abort()
			return
		}
		c.JSON(http.StatusCreated, newUser)
	}
}
