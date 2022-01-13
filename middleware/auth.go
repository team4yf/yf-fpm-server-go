package middleware

import (
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strings"
)

type BasicAuthConfig struct {
	Enable   bool
	Username string
	Pattern  []string
	Password string
}

func matchURL(url string, prefix []string) bool {
	for _, v := range prefix {
		matched, _ := filepath.Match(v, url)
		if matched {
			return true
		}
	}
	return false
}

func BasicAuth(basicAuth *BasicAuthConfig) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// skip options query
			if r.Method == "OPTIONS" {
				h.ServeHTTP(w, r)
				return
			}
			if !basicAuth.Enable || !matchURL(r.URL.String(), basicAuth.Pattern) {
				h.ServeHTTP(w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}
			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			if pair[0] != basicAuth.Username || pair[1] != basicAuth.Password {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
