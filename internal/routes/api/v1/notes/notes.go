package notes

import (
	"context"
	"fff-hedgedoc/internal/database"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/users"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetNotes(c *gin.Context) {
	_, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*12)
	defer cancel()

	notes, err := database.GetAllNotes(ctx)
	if err != nil {
		logging.Logger.Warnf("Error loading notes from database: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, note := range notes {
		logging.Logger.Infof("note: %s", note)
	}

	c.Status(http.StatusOK)
}

func SetRoutes(r *gin.RouterGroup) {
	notesR := r.Group("/notes")
	notesR.GET("", GetNotes)
}
