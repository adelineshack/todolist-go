package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	core_http_request "github.com/adelineshack/todolist-go/internal/core/transport/http/request"
	core_http_response "github.com/adelineshack/todolist-go/internal/core/transport/http/response"
)

type GetTasksReponse = []TasksDTOResponse

func (h *TasksHTTPHandler) GetTasks(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userId, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId/limit/offset query params")
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksReponse(tasksDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)

}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {

	const (
		usetIDQueryParamKey = "user_id"
		limitQueryParamkey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, usetIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamkey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
