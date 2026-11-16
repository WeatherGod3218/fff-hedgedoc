package users

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/models"
	"fff-hedgedoc/internal/redis"
	"fff-hedgedoc/internal/users"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	// _, err := users.GetCSHAuth(c)
	// if err != nil {
	// 	c.Status(http.StatusUnauthorized)
	// 	return
	// }

	var req models.GetUsersInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error binding JSON: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	redis_users, err := redis.GetUsersSection(ctx, req.Prefix, req.Cursor)
	if err != nil {
		logging.Logger.Warnf("Error getting users from redis: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	users := make([]*models.GetUsersPiece, len(redis_users))
	for index, user := range redis_users {
		users[index] = &models.GetUsersPiece{
			Username: user.Username,
			UUID:     user.UUID,
			Selected: false,
		}
	}

	cursor := ""
	if len(users) > 0 {
		cursor = users[len(users)-1].Username
	}

	c.JSON(http.StatusOK, models.GetUsersOutput{
		Users:  users,
		Cursor: cursor,
	})
}

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
	usersR := r.Group("/users")
	usersR.Handle("QUERY", "", GetUsers)
	usersR.GET("/info", GetUserInfo)
}
