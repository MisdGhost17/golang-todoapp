package auth_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type RegisterRequest struct {
	FullName    string  `json:"full_name" binding:"required" example:"Ivan Petrov"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
	Login       string  `json:"login" binding:"required" example:"auth_test"`
	// Password должен содержать минимум 15 символов и максимум 72 байта UTF-8.
	Password string `json:"password" binding:"required" example:"learning-go-auth-2026"`
}

type RegisterResponse struct {
	ID          int     `json:"id" example:"42"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"Ivan Petrov"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
}

// Register godoc
// @Summary Регистрация
// @Description Создаёт профиль и данные входа атомарно. Логин нормализуется: пробелы по краям удаляются, регистр приводится к нижнему. Пароль: минимум 15 символов, максимум 72 байта. Токены не выдаются.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Тело запроса"
// @Success 201 {object} RegisterResponse "Профиль создан"
// @Failure 400 {object} core_http_response.ErrorResponse "Некорректные данные"
// @Failure 409 {object} core_http_response.ErrorResponse "Логин уже занят"
// @Failure 500 {object} core_http_response.ErrorResponse "Внутренняя ошибка"
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	rw.Header().Set("Content-Type", "application/json")
	r.Body = http.MaxBytesReader(rw, r.Body, 16*1024)
	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode registration request")
		return
	}

	user, err := h.authService.Register(
		ctx,
		request.FullName,
		request.PhoneNumber,
		request.Login,
		request.Password,
	)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	responseHandler.JSONResponse(RegisterResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}, http.StatusCreated)
}
