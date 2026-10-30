package main

import (
	"net/http/httputil"
	"net/url"
	"os"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	hedgedocURI := os.Getenv("HEDGEDOC_SERVICE_URI")
	target, err := url.Parse(hedgedocURI)
	if err != nil {
		Logger.Fatalf("Error trying to intitalize reverse proxy %s", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.SetXForwarded()
		},
	}

	hostUrl := os.Getenv("SERVER_HOST")
	auth, err := csh_auth.Init(
		os.Getenv("AUTH_OIDC_ID"),
		os.Getenv("AUTH_OIDC_SECRET"),
		hostUrl,
		hostUrl+"/auth/login",
		hostUrl+"/auth/callback",
		[]string{"profile", "email", "groups"},
	)
	if err != nil {
		Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing csh-auth")
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/auth/login", auth.HandleLogin)
	router.GET("/auth/callback", auth.HandleCallback)
	router.GET("/auth/logout", auth.HandleLogout)

	cookieAuth := auth.CookieMiddleware()
	router.Use(func(c *gin.Context) {
		cookieAuth(c)
		if c.Writer.Written() || (c.Writer.Status() >= 300 && c.Writer.Status() < 400) {
			c.Abort()
		}
	})

	router.NoRoute(func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	router.Run(":8080")
}
