package handlers

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"backend/internal/logger"
	"github.com/sirupsen/logrus"
)

func obfuscateIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return host
	}
	if ip4 := ip.To4(); ip4 != nil {
		return fmt.Sprintf("%d.%d.%d.xxx", ip4[0], ip4[1], ip4[2])
	}
	return fmt.Sprintf("%x:%x:%x:%x:xxxx:xxxx:xxxx:xxxx", ip[0], ip[1], ip[2], ip[3])
}

func obfuscateUA(ua string) string {
	if len(ua) > 60 {
		return ua[:60] + "..."
	}
	return ua
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
				"ip":      obfuscateIP(r.RemoteAddr),
				"agent":   obfuscateUA(r.UserAgent()),
				"latency": time.Since(start).String(),
			}).Info("request")
		})
	}
}
