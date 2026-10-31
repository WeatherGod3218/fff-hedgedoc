package main

import (
	"bytes"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/users"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const COOKIE_ORIGINAL_REFERENCE = "original_reference"

var old = []byte(`<ul class="nav masthead-nav">`)

var inject = []byte(`
<ul class="nav masthead-nav"> 
<li class="ui-home">
	<a href="/csh/notes">CSH</a> 
</li>`)

func isSafePath(p string) bool { //had AI look over for bugs and it recommended this? alr
	return strings.HasPrefix(p, "/") &&
		!strings.HasPrefix(p, "//") &&
		!strings.HasPrefix(p, "/\\") &&
		!strings.HasPrefix(p, "/auth/")
}

func main() {
	hedgedocURI := os.Getenv("HEDGEDOC_SERVICE_URI")
	target, err := url.Parse(hedgedocURI)
	if err != nil {
		logging.Logger.Fatalf("Error trying to intitalize reverse proxy %s", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.SetXForwarded()
			scheme := "http"
			if r.In.Header.Get("X-Forwarded-Proto") == "https" || r.In.TLS != nil {
				scheme = "https"
			}
			r.Out.Header.Set("X-Forwarded-Proto", scheme)
			r.Out.Header.Set("Accept-Encoding", "identity")

			if r.In.URL.Path == "/auth/oauth2" || r.In.URL.Path == "/auth/oauth2/" {
				if cookie, err := r.In.Cookie(COOKIE_ORIGINAL_REFERENCE); err == nil {
					if dest, err := url.QueryUnescape(cookie.Value); err == nil && isSafePath(dest) {
						r.Out.Header.Set("Referer", os.Getenv("SERVER_HOST")+dest)
					}
				}
			}
		},

		ModifyResponse: func(resp *http.Response) error {
			if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
				return nil
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			body = bytes.Replace(body, old, inject, 1)

			resp.Body = io.NopCloser(bytes.NewReader(body))
			resp.ContentLength = int64(len(body))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
			return nil
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
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing csh-auth")
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(cors.Default())

	secureCookie := strings.HasPrefix(os.Getenv("SERVER_HOST"), "https://")

	router.GET("/auth/login", func(ctx *gin.Context) {
		if dest := ctx.Request.URL.Query().Get("referer"); isSafePath(dest) {
			ctx.SetCookie(COOKIE_ORIGINAL_REFERENCE, dest, 400, "/", "", secureCookie, true)
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

		if valid := users.IsFFF(user); valid != true {
			ctx.HTML(http.StatusForbidden, "forbidden.tmpl", gin.H{
				"User": user.Username,
			})
			return
		}

		if ctx.Request.URL.Path == "/auth/oauth2/callback" {
			ctx.SetCookie(COOKIE_ORIGINAL_REFERENCE, "", -1, "/", "", secureCookie, true)
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	router.Run(":8080")
}
