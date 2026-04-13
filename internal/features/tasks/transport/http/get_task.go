package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	core_http_request "github.com/adelineshack/todolist-go/internal/core/transport/http/request"
	core_http_response "github.com/adelineshack/todolist-go/internal/core/transport/http/response"
)

type GetTaskReponse TasksDTOResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskReponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)
}
