package handlers

import (
	"net/http"
	"runtime/debug"
	"time"

	"backend/internal/logger"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.WithFields(logrus.Fields{
						"ip":       logger.ObfuscateIP(r.RemoteAddr),
						"path":     r.URL.Path,
						"method":   r.Method,
						"panic":    rec,
						"stack":    string(debug.Stack()),
					}).Error("panic recovered")
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			log.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"status":  rw.status,
				"ip":      logger.ObfuscateIP(r.RemoteAddr),
				"agent":   logger.ObfuscateUA(r.UserAgent()),
				"latency": time.Since(start).String(),
			}).Info("request")
		})
	}
}
