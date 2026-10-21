package users_transport_http

import (
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
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
