package groups

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

func GetGroups(c *gin.Context) {
	// _, err := users.GetCSHAuth(c)
	// if err != nil {
	// 	c.Status(http.StatusUnauthorized)
	// 	return
	// }

	var req models.GetGroupsInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error binding JSON: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	redis_groups, err := redis.GetGroupsSection(ctx, req.Prefix, req.Cursor)
	if err != nil {
		logging.Logger.Warnf("Error getting users from redis: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	groups := make([]*models.GetGroupsPiece, len(redis_groups))
	for index, user := range redis_groups {
		groups[index] = &models.GetGroupsPiece{
			GroupName: user.GroupName,
			UUID:      user.UUID,
			Selected:  false,
		}
	}

	cursor := ""
	if len(groups) > 0 {
		cursor = groups[len(groups)-1].GroupName
	}

	c.JSON(http.StatusOK, models.GetGroupsOutput{
		Groups: groups,
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
	groupsR := r.Group("/groups")
	groupsR.Handle("QUERY", "", GetGroups)

}
