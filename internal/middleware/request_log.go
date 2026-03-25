package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ulbithebest/BE-pendaftaran/internal/applog"
	"github.com/ulbithebest/BE-pendaftaran/internal/auth"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rec *statusRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *statusRecorder) Write(body []byte) (int, error) {
	if rec.statusCode == 0 {
		rec.statusCode = http.StatusOK
	}

	written, err := rec.ResponseWriter.Write(body)
	rec.bytesWritten += written
	return written, err
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

		fullPath := r.URL.Path
		if r.URL.RawQuery != "" {
			fullPath = fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
		}

		payload := extractPayloadFromAuthorization(r.Header.Get("Authorization"))
		durationMs := time.Since(start).Milliseconds()

		entry := applog.Entry{
			Timestamp:     time.Now(),
			Level:         level,
			Message:       fmt.Sprintf("%s %s -> %d (%d ms)", r.Method, fullPath, recorder.statusCode, durationMs),
			Method:        r.Method,
			Path:          r.URL.Path,
			Query:         r.URL.RawQuery,
			FullPath:      fullPath,
			StatusCode:    recorder.statusCode,
			DurationMs:    durationMs,
			RemoteIP:      getRequestIP(r),
			UserAgent:     r.UserAgent(),
			Host:          r.Host,
			Referer:       r.Referer(),
			ResponseBytes: recorder.bytesWritten,
			ContentLength: r.ContentLength,
		}

		if payload != nil {
			entry.UserID = payload.UserID.Hex()
			entry.UserNIM = payload.NIM
			entry.UserRole = payload.Role
		}

		applog.Add(entry)
	})
}

func extractPayloadFromAuthorization(authHeader string) *auth.PasetoPayload {
	if authHeader == "" {
		return nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}

	payload, err := auth.VerifyToken(parts[1])
	if err != nil {
		return nil
	}

	return payload
}

func getRequestIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if forwardedFor != "" {
		return forwardedFor
	}

	xRealIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xRealIP != "" {
		return xRealIP
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
