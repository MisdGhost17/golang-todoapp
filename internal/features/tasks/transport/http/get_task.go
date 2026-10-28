package tasks_transport_http

import (
	"net/http"

	core_auth "github.com/MisdGhost17/golang-todoapp/internal/core/auth"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"

	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получение задачи
// @Description Получение своей задачи по ID. Для чужой или отсутствующей задачи возвращается 404.
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse "Задача успешно найдена"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "authentication required")
		return
	}

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
