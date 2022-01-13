package middleware

import (
	"net/http"
)

type ServerAuthConfig struct {
	Enable bool                         `json:"enable"`
	Match  map[string]ServerMatchConfig `json:"match"`
}

type ServerMatchConfig struct {
	Keys    []string `json:"keys"`
	Header  string   `json:"header"`
	Pattern []string `json:"pattern"`
}

func ServerAuth(serverAuth *ServerAuthConfig) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// skip options query
			if r.Method == "OPTIONS" {
				h.ServeHTTP(w, r)
				return
			}
			if !serverAuth.Enable {
				h.ServeHTTP(w, r)
				return
			}
			for _, v := range serverAuth.Match {
				if matchURL(r.URL.String(), v.Pattern) {
					key := r.Header.Get(v.Header)
					if key == "" {
						http.Error(w, "Not authorized! Missing server key!", http.StatusUnauthorized)
						return
					}
					matched := false
					for _, k := range v.Keys {
						if k == key {
							matched = true
							break
						}
					}
					if !matched {
						http.Error(w, "Not authorized! Server key not matched", http.StatusUnauthorized)
						return
					}
					// matched, break out of loop
					break
				}
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
