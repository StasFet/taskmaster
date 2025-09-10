package internal

import (
	"net/http"
	"os"
	"strings"
	"taskmaster/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// verifies the jwt token provided by supabase auth
func ValidateToken(token string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
}

// check that all the claims are right.
func ValidateClaims(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// checking if token is expired or not
		expiry := claims["exp"].(float64)
		expiryTime := time.Unix(int64(expiry), 0)
		if time.Now().After(expiryTime) {
			return false, nil
		}

		// is token not valid yet,
		nbf := claims["nbf"].(float64)
		notBefore := time.Unix(int64(nbf), 0)
		if time.Now().Before(notBefore) {
			return false, nil
		}

		// check the audience is correct (supabase sets it to "authenticated" by defauly)
		aud := claims["aud"].(string)
		if aud != "authenticated" {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// middleware function to check the jwt with the request and to attach the extracted uuid to the context
func JWTValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		parsedToken, err := ValidateToken(authToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"message": "authentication failed: invalid token",
			})
			logger.GEN.Printf("Error validating token: %v\n", err)
			c.Abort()
			return
		}
		// check expiry etc
		if isValid, err := ValidateClaims(parsedToken); !isValid || err != nil {
			c.JSON(http.StatusForbidden, map[string]any{
				"message": "authentication failed: expired token",
			})
			logger.GEN.Printf("Error validating token claims. Token: %v\n", parsedToken)
			c.Abort()
			return
		}

		// set the uuid for the following functions to use
		claims, _ := parsedToken.Claims.(jwt.MapClaims)
		uuid := claims["sub"].(string)
		name := claims["name"].(string)
		email := claims["email"].(string)
		c.Set("validated_uuid", uuid)
		c.Set("validated_name", name)
		c.Set("validated_email", email)
		c.Next()
	}
}
