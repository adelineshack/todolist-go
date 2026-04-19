package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/adelineshack/todolist-go/internal/core/config"
	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	core_pgxpool "github.com/adelineshack/todolist-go/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/adelineshack/todolist-go/internal/core/transport/http/middleware"
	core_http_server "github.com/adelineshack/todolist-go/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/adelineshack/todolist-go/internal/features/statistics/repository"
	statistics_service "github.com/adelineshack/todolist-go/internal/features/statistics/service"
	statistics_transport_http "github.com/adelineshack/todolist-go/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/adelineshack/todolist-go/internal/features/tasks/repository/postgres"
	tasks_service "github.com/adelineshack/todolist-go/internal/features/tasks/service"
	tasks_transport_http "github.com/adelineshack/todolist-go/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/adelineshack/todolist-go/internal/features/users/repository/postgres"
	users_service "github.com/adelineshack/todolist-go/internal/features/users/service"
	users_transport_http "github.com/adelineshack/todolist-go/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Println("failer to init application logger:", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("Initializing postgres connection pool")

	pool, err := core_pgxpool.NewPool(
		ctx,
		core_pgxpool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("Initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewTasksService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("Initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routers()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
