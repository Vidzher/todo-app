package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	core_loger "todo-api/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_loger.Logger
}

func NewHTTPServer(config Config, logger *core_loger.Logger) *HTTPServer {
	return &HTTPServer{mux: http.NewServeMux(), config: config, log: logger}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Address,
		Handler: s.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("address", s.config.Address))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}

	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}

func (s *HTTPServer) RegisterAPIRouters(routers ...APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}
