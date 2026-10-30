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
			r.Out.Header.Set("X-Forwarded-Proto", "https")
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

	router.GET("/auth/login", func(c *gin.Context) {
		c.Request.URL.RawQuery = "referer=/auth/oauth2"
		auth.HandleLogin(c)
	})

	router.GET("/auth/callback", auth.HandleCallback)
	router.GET("/auth/logout", auth.HandleLogout)

	router.Use(auth.CookieMiddleware())

	router.NoRoute(func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	router.Run(":8080")
}
