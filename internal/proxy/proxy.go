package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const COOKIE_ORIGINAL_REFERENCE = "original_reference"

var old = []byte(`<ul class="nav masthead-nav">`)

var inject = []byte(`
<ul class="nav masthead-nav"> 
<li class="ui-home">
	<a href="/csh/notes">CSH</a> 
</li>`)

func IsSafePath(p string) bool { //had AI look over for bugs and it recommended this? alr
	return strings.HasPrefix(p, "/") &&
		!strings.HasPrefix(p, "//") &&
		!strings.HasPrefix(p, "/\\") &&
		!strings.HasPrefix(p, "/auth/")
}

func GetProxy() (*httputil.ReverseProxy, error) {
	hedgedocURI := os.Getenv("HEDGEDOC_SERVICE_URI")
	target, err := url.Parse(hedgedocURI)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.SetXForwarded()
			scheme := "http"
			if r.In.Header.Get("X-Forwarded-Proto") == "https" || r.In.TLS != nil {
				scheme = "https"
			}
			r.Out.Header.Set("X-Forwarded-Proto", scheme)
			r.Out.Header.Set("Accept-Encoding", "identity")
		},

		ModifyResponse: func(resp *http.Response) error {
			if strings.HasSuffix(resp.Request.URL.Path, "/auth/oauth2/callback") &&
				resp.StatusCode >= 300 && resp.StatusCode < 400 {
				if c, err := resp.Request.Cookie(COOKIE_ORIGINAL_REFERENCE); err == nil {
					if dest, err := url.QueryUnescape(c.Value); err == nil && IsSafePath(dest) {
						resp.Header.Set("Location", dest)
					}
				}
			}

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
	}, nil
}
