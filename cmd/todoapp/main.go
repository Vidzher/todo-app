package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	core_loger "todo-api/internal/core/logger"
	core_postgres_pool "todo-api/internal/core/repository/postgres/pool"
	core_http_middleware "todo-api/internal/core/transport/http/middleware"
	core_http_server "todo-api/internal/core/transport/http/server"
	tasks_postgres_repository "todo-api/internal/features/tasks/repository/postgres"
	tasks_service "todo-api/internal/features/tasks/service"
	tasks_transport_http "todo-api/internal/features/tasks/transport/http"
	users_postgres_repository "todo-api/internal/features/users/repository/postgres"
	users_service "todo-api/internal/features/users/service"
	users_transport_http "todo-api/internal/features/users/transport/http"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_loger.NewLogger(core_loger.NewConfigMust())
	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to start HTTP server", zap.Error(err))
	}
}
