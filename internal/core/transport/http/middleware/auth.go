package core_http_middleware

import (
	"net/http"
	"strings"

	core_auth "github.com/MisdGhost17/golang-todoapp/internal/core/auth"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type TokenVerifier interface {
	Verify(tokenString string) (int, error)
}

func RequireAuth(verifier TokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reject := func() {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("WWW-Authenticate", "Bearer")

				log := core_logger.FromContext(r.Context())
				responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
				responseHandler.ErrorResponse(
					core_errors.ErrUnauthorized,
					"valid bearer token required",
				)
			}

			headers := r.Header.Values("Authorization")
			if len(headers) != 1 {
				reject()
				return
			}

			parts := strings.Fields(headers[0])
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				reject()
				return
			}

			userID, err := verifier.Verify(parts[1])
			if err != nil || userID <= 0 {
				reject()
				return
			}

			ctx := core_auth.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
