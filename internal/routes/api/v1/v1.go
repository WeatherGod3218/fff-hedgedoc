package v1

import (
	"fff-hedgedoc/internal/routes/api/v1/notes"
	"fff-hedgedoc/internal/routes/api/v1/users"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.RouterGroup) {
	v1R := r.Group("/v1")
	notes.SetRoutes(v1R)
	users.SetRoutes(v1R)
}
