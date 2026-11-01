package routes

import (
	"fff-hedgedoc/internal/routes/api"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.RouterGroup) {
	cshR := r.Group("/csh")
	api.SetRoutes(cshR)
}
