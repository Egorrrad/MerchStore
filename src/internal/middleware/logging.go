package middleware

import (
	"MerchStore/src/internal/logger"
	"bytes"
	"io"
	"net/http"
	"time"
)

const maxBodySize = 1024

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var requestBody string
		if r.Body != nil {
			limitedReader := io.LimitReader(r.Body, maxBodySize+1)
			bodyBytes, _ := io.ReadAll(limitedReader)

			if len(bodyBytes) > maxBodySize {
				requestBody = "[too large]"
			} else {
				requestBody = string(bodyBytes)
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		logger.Logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"user-agent", r.UserAgent(),
			"body", requestBody,
		).Info("Request started")

		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(lrw, r)

		logger.Logger.With(
			"status", lrw.status,
			"duration", time.Since(start),
		).Info("Request completed")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}
