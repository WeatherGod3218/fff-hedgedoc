package main

import (
	"embed"
	"fff-hedgedoc/internal/logging"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var distFS embed.FS

var icon []byte
var index []byte

func ServeFrontend(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", index)
}

func ServeIcon(c *gin.Context) {
	c.Data(http.StatusOK, "image/svg", icon)
}

func LoadFrontend(r *gin.RouterGroup) {
	dist, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		logging.Logger.Fatalf("failed to initialize dist for the app %s", err)
	}

	assets, err := fs.Sub(dist, "assets")
	if err != nil {
		logging.Logger.Fatalf("failed to initialize assets for the app %s", err)
	}

	index, err = fs.ReadFile(dist, "index.html")
	if err != nil {
		logging.Logger.Fatalf("failed to load index %s", err)
	}

	icon, err = fs.ReadFile(dist, "favicon.svg")
	if err != nil {
		logging.Logger.Fatal("No favicon?")
	}

	r.Use(cors.Default())
	r.StaticFS("/assets", http.FS(assets))

	r.GET("/favicon.ico", func(ctx *gin.Context) {
		ServeIcon(ctx)
	})

	r.GET("/favicon.svg", func(ctx *gin.Context) { // just do both and make my life easy
		ServeIcon(ctx)
	})

	r.GET("/favicon.png", func(ctx *gin.Context) { // theres a third now ig.
		ServeIcon(ctx)
	})

	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
