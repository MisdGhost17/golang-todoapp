package auth_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required" example:"auth_test"`
	Password string `json:"password" binding:"required" example:"learning-go-auth-2026"`
}

type LoginResponse struct {
	// AccessToken — JWT для заголовка Authorization: Bearer <access_token>.
	AccessToken string `json:"access_token"`
	// RefreshToken — непрозрачный одноразовый токен: 43 символа Base64URL без padding.
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// Login godoc
// @Summary Вход
// @Description Проверяет логин и пароль, создаёт новую сессию и выдаёт пару токенов. Access JWT действует 15 минут, refresh-сессия — 30 дней. Токены передаются в JSON, cookie не устанавливается. Для неизвестного логина и неверного пароля возвращается одинаковая ошибка.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Тело запроса"
// @Success 200 {object} LoginResponse "Пара токенов"
// @Header 200 {string} Cache-Control "no-store"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный JSON"
// @Failure 401 {object} core_http_response.ErrorResponse "Неверный логин или пароль"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка"
// @Router /auth/login [post]
func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Cache-Control", "no-store")
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate login request")
		return
	}

	tokens, err := h.authService.Login(
		ctx,
		request.Login,
		request.Password,
	)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	responseHandler.JSONResponse(LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
	}, http.StatusOK)
}
