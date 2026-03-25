package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ulbithebest/BE-pendaftaran/internal/applog"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		level := "info"
		switch {
		case recorder.statusCode >= 500:
			level = "error"
		case recorder.statusCode >= 400:
			level = "warn"
		}

		applog.Add(applog.Entry{
			Timestamp:  time.Now(),
			Level:      level,
			Message:    fmt.Sprintf("%s %s -> %d", r.Method, r.URL.Path, recorder.statusCode),
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: recorder.statusCode,
			DurationMs: time.Since(start).Milliseconds(),
			RemoteIP:   r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		})
	})
}
