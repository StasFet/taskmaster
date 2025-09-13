package security

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"taskmaster/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
	Exp   int64  `json:"exp"`
	jwt.RegisteredClaims
}

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
		expiry, ok := claims["exp"].(float64)
		if !ok {
			return false, nil
		}
		expiryTime := time.Unix(int64(expiry), 0)
		if time.Now().After(expiryTime) {
			return false, nil
		}

		// iat is issued at; check if the key is issued now
		iat, ok := claims["iat"].(float64)
		if !ok {
			return false, nil
		}
		issuedAt := time.Unix(int64(iat), 0)
		if time.Now().Before(issuedAt) {
			return false, nil
		}

		// check the audience is correct (supabase sets it to "authenticated" by defauly)
		aud, ok := claims["aud"].(string)
		if aud != "authenticated" || !ok {
			return false, nil
		}

		// check if the issuer is valid
		iss, ok := claims["iss"].(string)
		if iss != "https://onqmqxugejuudbvxuzyq.supabase.co/auth/v1" || !ok {
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
		//logger.GEN.Printf("Recieved authorisation, here is the token: %v\n", authToken)
		if len(authToken) == 0 {
			logger.API.Printf("Error processing JWT: token had length 0")
			c.Abort()
			return
		}
		parsedToken, err := ValidateToken(authToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"message": "authentication failed: invalid token",
			})
			logger.GEN.Printf("Error validating token: %v\n", err)
			c.Abort()
			return
		}

		logToken, err := json.Marshal(parsedToken)
		if err == nil {
			logger.API.Printf("Extracted logToken: %s", logToken)
		}

		// check expiry etc
		if isValid, err := ValidateClaims(parsedToken); !isValid || err != nil {
			c.JSON(http.StatusForbidden, map[string]any{
				"message": "authentication failed: expired",
			})
			logger.GEN.Printf("Error validating token claims. Token: %v\n", parsedToken)
			c.Abort()
			return
		}

		// set the uuid for the following functions to use
		claims, _ := parsedToken.Claims.(jwt.MapClaims)
		uuid := claims["sub"].(string)
		email := claims["email"].(string)
		c.Set("validated_uuid", uuid)
		c.Set("validated_email", email)

		// name is in user_metadata of claims, so we need to extract it a bit weird
		name := ""
		if user_metadata, ok := claims["user_metadata"].(map[string]any); ok {
			if n, ok := user_metadata["name"].(string); ok {
				name = n
			} else {
				logger.API.Printf("Error extracting name from \"user_metadata\" of jwt key")
			}
		}
		c.Set("validated_name", name)
		logger.GEN.Printf("Token validated successfully! email: %s\n", email)
		c.Next()
	}
}
