package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"www-api/internal/sso"

	"github.com/gin-gonic/gin"
)

// Authenticate middleware function to authenticate all request based on passed token
func Authenticate(c *gin.Context) {
	//fetch authorization header
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		header = c.Request.Header.Get("authorization")
	}

	//check if header is not passed
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "missing Authorization header",
		})
		return
	}

	auth := strings.Split(header, " ")
	//check for bearer token
	if auth[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid Authorization header format",
		})
		return
	}

	//fetch deployment type from gin context
	environement, ok := c.Get("deployment")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "type of deployment not available in config",
		})
		return
	}

	//verify token based on deployment type
	isValid, err := sso.VerifyToken(auth[1], environement.(string))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "err while validating token",
		})
		return
	}

	//check if token is valid
	if !isValid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}

	c.Next()

}
