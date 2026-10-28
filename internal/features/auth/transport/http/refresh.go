package auth_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type RefreshRequest struct {
	// RefreshToken — непрозрачный одноразовый токен: 43 символа Base64URL без padding.
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	// AccessToken — JWT для заголовка Authorization: Bearer <access_token>.
	AccessToken string `json:"access_token"`
	// RefreshToken — непрозрачный одноразовый токен: 43 символа Base64URL без padding.
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// Refresh godoc
// @Summary Обновление токенов
// @Description Обменивает refresh-токен на новую пару. Access JWT не требуется. Старый refresh-токен становится использованным; повторное использование отзывает всю сессию. Срок сессии не продлевается. Не отправляйте параллельные запросы и не повторяйте запрос автоматически при потере ответа.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Тело запроса"
// @Success 200 {object} RefreshResponse "Новая пара токенов"
// @Header 200 {string} Cache-Control "no-store"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный JSON"
// @Failure 401 {object} core_http_response.ErrorResponse "Токен недействителен, истёк, использован или сессия отозвана"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка"
// @Router /auth/refresh [post]
func (h *AuthHTTPHandler) Refresh(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Cache-Control", "no-store")
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode refresh request")
		return
	}

	tokens, err := h.authService.Refresh(ctx, request.RefreshToken)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh tokens")
		return
	}

	responseHandler.JSONResponse(RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
	}, http.StatusOK)
}
