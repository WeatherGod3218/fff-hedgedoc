package access

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/models"
	"fff-hedgedoc/internal/redis"
	"fff-hedgedoc/internal/users"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAccessModifier(c *gin.Context) {
	_, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.ModifyAccessInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error binding req: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	fullName := redis.GROUP_MEMBERS_PREFIX + req.IdentifierUUID
	members, err := redis.GetAllFromKeyPrefix(ctx, fullName, -1)
	if err != nil {
		logging.Logger.Warnf("Error getting values from key prefix: %s", err)
	}

	for key, id := range members {
		logging.Logger.Infof("Member: %s, Id: %s", strings.TrimPrefix(key, fullName), id)
	}

	logging.Logger.Info("Should have printed all")
	c.Status(http.StatusOK)
}

func DeleteAccessModifier(c *gin.Context) {
	_, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

}

func SetRoutes(r *gin.RouterGroup) {
	accessR := r.Group("/access")
	accessR.POST("", AddAccessModifier)
	accessR.DELETE("", AddAccessModifier)
}
