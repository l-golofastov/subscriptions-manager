package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewLoggingMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		requestID := GetRequestID(r.Context())

		info := fmt.Sprintf("method=%s path=%s status=%d duration=%s request_id=%s",
			r.Method, r.URL.Path, rw.statusCode, duration, requestID)

		log.Info(info)
	})
}
