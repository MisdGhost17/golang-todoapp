package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/MisdGhost17/golang-todoapp/internal/core/auth"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Получение списка задач
// @Description Получение задач текущего пользователя с пагинацией. Параметр user_id игнорируется. Без limit возвращаются все свои задачи.
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Максимальное количество задач" minimum(0)
// @Param offset query int false "Количество пропускаемых задач" minimum(0)
// @Success 200 {array} TaskDTOResponse "Список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrUnauthorized,
			"authentication required",
		)

		return
	}

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get pagination query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
