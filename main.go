package main

import (
	"fff-hedgedoc/internal/database"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/proxy"
	"fff-hedgedoc/internal/redis"
	"fff-hedgedoc/internal/routes"
	"fff-hedgedoc/internal/users"
	"net/http"
	"os"
	"strings"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func main() {
	pxy, err := proxy.GetProxy()
	if err != nil {
		logging.Logger.Fatalf("Failed to load proxy with error: %s", err)
	}

	if err := redis.InitRedis(); err != nil {
		logging.Logger.Warnf("Failed to load redis client with error: %s, proceeding without", err)
	}

	if err := database.InitDatabase(); err != nil {
		logging.Logger.Fatalf("Failed to load postgres database with error: %s", err)
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
		logging.Logger.Fatalf("Failed to load csh-auth with error: %s", err)
	}

	router := gin.Default()

	secureCookie := strings.HasPrefix(hostUrl, "https://")

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

	LoadFrontend(router.Group(""))
	routes.SetRoutes(router.Group(""))

	if os.Getenv("DEV_MODE") == "true" {
		logging.Logger.Info("RUNNING IN DEV MODE!")
		router.NoRoute(createViteProxy())
	} else {
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

			if strings.HasPrefix(ctx.Request.URL.Path, "/csh") {
				ServeFrontend(ctx)
				return
			}

			pxy.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	router.Run(":8080")
}
