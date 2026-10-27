package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Получение списка задач
// @Description Получение списка задач с опциональной фильтрацией по автору и пагинацией.
// @Tags tasks
// @Produce json
// @Param user_id query int false "ID автора задач"
// @Param limit query int false "Максимальное количество задач" minimum(0)
// @Param offset query int false "Количество пропускаемых задач" minimum(0)
// @Success 200 {array} TaskDTOResponse "Список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
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

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDParamKey      = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
