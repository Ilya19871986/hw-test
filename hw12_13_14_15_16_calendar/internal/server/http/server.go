package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
)

type Server struct {
	server *http.Server
	logger *logger.Logger
	app    *app.App
}

func NewServer(host, port string, app *app.App, logger *logger.Logger) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: loggingMiddleware(mux, logger),
	}

	mux.HandleFunc("/hello", helloHandler)

	return &Server{
		server: server,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server on " + s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
