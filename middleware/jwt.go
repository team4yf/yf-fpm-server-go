package middleware

import (
	"net/http"
	"strings"

	"github.com/team4yf/fpm-go-pkg/utils"
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
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Not authorized! Token missing!", http.StatusUnauthorized)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			if ok, _ := utils.CheckToken(token); !ok {
				http.Error(w, "Not authorized! Token is invalid", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
