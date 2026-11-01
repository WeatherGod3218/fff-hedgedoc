package api

import (
	v1 "fff-hedgedoc/internal/routes/api/v1"

	"github.com/gin-gonic/gin"
)

func CookieToHeaderAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("Auth")
		if err == nil && cookie != "" {
			c.Request.Header.Set("Authorization", "Bearer "+cookie)
		}

		c.Next()
	}
}

func SetRoutes(r *gin.RouterGroup) {
	apiR := r.Group("/api", CookieToHeaderAuth())

	v1.SetRoutes(apiR)
}
