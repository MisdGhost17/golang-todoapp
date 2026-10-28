package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_http_server "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/server"
)

type AuthService interface {
	Register(
		ctx context.Context,
		fullName string,
		phoneNumber *string,
		login string,
		password string,
	) (domain.User, error)

	Login(
		ctx context.Context,
		login string,
		password string,
	) (domain.TokenPair, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (domain.TokenPair, error)

	Logout(
		ctx context.Context,
		refreshToken string,
	) error
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}
