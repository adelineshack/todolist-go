package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	core_http_response "github.com/adelineshack/todolist-go/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rosponseHandler := core_http_response.NewHttpResponseHandler(log, rw)

}
