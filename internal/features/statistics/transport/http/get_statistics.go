package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"20"`
	TasksCompleted             int      `json:"tasks_completed" example:"15"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"75"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1h30m0s"`
}

// GetStatistics godoc
// @Summary Получение статистики задач
// @Description Получение статистики задач с опциональными фильтрами по пользователю и периоду.
// @Tags statistics
// @Produce json
// @Param user_id query int false "ID пользователя"
// @Param from query string false "Начало периода в формате YYYY-MM-DD" Format(date)
// @Param to query string false "Конец периода в формате YYYY-MM-DD" Format(date)
// @Success 200 {object} GetStatisticsResponse "Статистика задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFRomToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	statisticsDomain, err := h.statisticsService.GetTasks(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"get statistics from service",
		)

		return
	}

	response := toDTOFromDomain(statisticsDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageComplitionTime != nil {
		duartion := statistics.TasksAverageComplitionTime.String()
		avgTime = &duartion
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFRomToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
