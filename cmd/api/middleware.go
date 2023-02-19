package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthErrorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Error: "No Token Provide"})
			return
		}
		tokenString = tokenString[(len("Bearer ")):]

		token, err := ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Error: "Invalid Token"})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, AuthErrorResponse{Error: "Bad Request"})
			return
		}
		if token.Valid {
			// claims := token.Claims.(*Claims)
			// fmt.Println(claims.Email)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Error: "Invalid token"})
		}
	}
}
