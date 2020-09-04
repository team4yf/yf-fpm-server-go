package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/team4yf/yf-fpm-server-go/pkg/log"
)

//Recover recover the panic
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if !strings.HasPrefix(r.URL.String(), "/api") && r.Method == "POST" {
					log.Errorf("URL: %s, METHOD: %s, Error: %+v\n", r.URL.String(), r.Method, err)
					http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
					return
				}
				//
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
