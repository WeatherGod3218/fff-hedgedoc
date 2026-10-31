package main

import (
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/proxy"
	"fff-hedgedoc/internal/users"
	"net/http"
	"os"
	"strings"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	pxy, err := proxy.GetProxy()

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
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing csh-auth")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(cors.Default())

	secureCookie := strings.HasPrefix(os.Getenv("SERVER_HOST"), "https://")

	router.GET("/auth/login", func(ctx *gin.Context) {
		if dest := ctx.Request.URL.Query().Get("referer"); proxy.IsSafePath(dest) {
			ctx.SetCookie(proxy.COOKIE_ORIGINAL_REFERENCE, dest, 400, "/", "", secureCookie, true)
		}
		ctx.Request.URL.RawQuery = "referer=/auth/oauth2"
		auth.HandleLogin(ctx)
	})

	router.GET("/auth/callback", auth.HandleCallback)
	router.GET("/auth/logout", auth.HandleLogout)

	router.Use(auth.CookieMiddleware())
	router.NoRoute(func(ctx *gin.Context) {
		user, err := users.GetCSHAuth(ctx)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		if !users.IsFFF(user) {
			ctx.HTML(http.StatusForbidden, "forbidden.tmpl", gin.H{
				"User": user.Username,
			})
			return
		}

		pxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	router.Run(":8080")
}
