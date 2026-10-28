package auth_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type LogoutRequest struct {
	// RefreshToken — непрозрачный одноразовый токен: 43 символа Base64URL без padding.
	RefreshToken string `json:"refresh_token"`
}

// Logout godoc
// @Summary Выход из сессии
// @Description Отзывает сессию по refresh-токену, в том числе использованному. Access JWT не требуется. Повторный выход, неизвестный или некорректный токен возвращают 204. Уже выданный access JWT действует до истечения срока. Клиент должен удалить оба токена.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Тело запроса"
// @Success 204 "Сессия отозвана или токен не найден; тело отсутствует"
// @Header 204 {string} Cache-Control "no-store"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректный JSON"
// @Failure 500 {object} core_http_response.ErrorResponse "Не удалось отозвать сессию"
// @Router /auth/logout [post]
func (h *AuthHTTPHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	rw.Header().Set("Cache-Control", "no-store")
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)

	var request LogoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rw.Header().Set("Content-Type", "application/json")
		responseHandler.ErrorResponse(err, "failed to decode logout request")
		return
	}

	if err := h.authService.Logout(ctx, request.RefreshToken); err != nil {
		rw.Header().Set("Content-Type", "application/json")
		responseHandler.ErrorResponse(err, "failed to logout")
		return
	}

	responseHandler.NoContentResponse()
}
