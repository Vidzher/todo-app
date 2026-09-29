package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	core_loger "todo-api/internal/core/logger"
	core_http_middleware "todo-api/internal/core/transport/http/middleware"
	core_http_server "todo-api/internal/core/transport/http/server"
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

	logger.Debug("starting todo-api")

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to start HTTP server", zap.Error(err))
	}
}
