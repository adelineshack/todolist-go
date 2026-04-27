package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adelineshack/todolist-go/internal/core/domain"
	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	core_http_request "github.com/adelineshack/todolist-go/internal/core/transport/http/request"
	core_http_response "github.com/adelineshack/todolist-go/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated              int
	TasksCompleted            int
	TasksCompletedRate        *float64
	TasksAverageCompletedTime *string
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userId, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId/from/to query params")
		return
	}

	tasksDomains, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(tasksDomains)

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string

	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:              statistics.TasksCreated,
		TasksCompleted:            statistics.TasksCompleted,
		TasksCompletedRate:        statistics.TasksCompletedRate,
		TasksAverageCompletedTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		usetIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, usetIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}

	to, err := core_http_request.GetDateQUeryParam(r, toQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	from, err := core_http_request.GetDateQUeryParam(r, fromQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, to, from, nil
}
