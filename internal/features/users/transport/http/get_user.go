package users_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser      godoc
// @Summary     Получение пользователя
// @Description Получение конкретного пользователя по его ID
// @Description В текущей конфигурации маршрут доступен без JWT.
// @Tags        users
// @Produce     json
// @Param       id path int true "ID пользователя"
// @Success     200 {object} GetUserResponse "Пользователь успешно найден"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 	404 {object} core_http_response.ErrorResponse "User Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 	 	/users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id path value")

		return
	}

	userDomain, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return
	}

	response := GetUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
