package internalhttp

import (
	"net/http"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		logger.Infof("%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			lw.statusCode,
			duration.Milliseconds(),
			r.UserAgent(),
		)
	})
}
