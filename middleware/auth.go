package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthConfig struct {
	Enable   bool
	Username string
	Password string
}

func BasicAuth(basicAuth *BasicAuthConfig) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !basicAuth.Enable || !strings.HasPrefix(r.URL.String(), "/biz") {
				h.ServeHTTP(w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}
			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				http.Error(w, err.Error(), 401)
				return
			}
			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}

			if pair[0] != basicAuth.Username || pair[1] != basicAuth.Password {
				http.Error(w, "Not authorized", 401)
				return
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
