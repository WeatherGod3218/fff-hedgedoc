package users

import (
	"fff-hedgedoc/internal/models"
	"fff-hedgedoc/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, models.UserInfo{
		UUID:     user.Uuid,
		Username: user.Username,
		FullName: user.FullName,
	})
}
func SetRoutes(r *gin.RouterGroup) {
	notesR := r.Group("/user")
	notesR.GET("", GetUserInfo)
}
