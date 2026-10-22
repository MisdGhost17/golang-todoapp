package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
)

func Dummy(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			log.Debug(fmt.Sprintf("-> before: %s", s))

			next.ServeHTTP(rw, r)

			log.Debug(fmt.Sprintf("<- after: %s", s))
		})
	}
}
