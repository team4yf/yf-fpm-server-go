package fpm

import (
	"net/http"

	"github.com/team4yf/yf-fpm-server-go/pkg/log"
)

//RecoverMiddleware recover the panic
func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Infof("URL: %s, METHOD: %s, Error: %+v\n", r.URL.String(), r.Method, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
