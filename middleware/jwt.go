package middleware

import (
	"net/http"

	"github.com/team4yf/fpm-go-pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/ctx"
)

type JwtAuthConfig struct {
	Enable  bool
	Pattern []string
}

func JwtAuth(jwtAuth *JwtAuthConfig) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// skip options query
			if r.Method == "OPTIONS" {
				h.ServeHTTP(w, r)
				return
			}
			if !jwtAuth.Enable || !matchURL(r.URL.String(), jwtAuth.Pattern) {
				h.ServeHTTP(w, r)
				return
			}
			token := ctx.WrapCtx(w, r).GetToken()
			if token == "" {
				http.Error(w, "Not authorized! Token missing!", http.StatusUnauthorized)
				return
			}
			if ok, _ := utils.CheckToken(token); !ok {
				http.Error(w, "Not authorized! Token is invalid", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
