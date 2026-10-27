package tasks_transport_http

import (
	"net/http"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Купить продукты"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Молоко и хлеб"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"7"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary Создание задачи
// @Description Создание задачи для существующего пользователя.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Данные новой задачи"
// @Success 201 {object} CreateTaskResponse "Задача успешно создана"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "Author User Not Found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewTaskUnitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
