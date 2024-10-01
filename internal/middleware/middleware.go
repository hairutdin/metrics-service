package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := wrapResponseWriter(w)
			next.ServeHTTP(wrappedWriter, r)
			duration := time.Since(start)

			logger.WithFields(logrus.Fields{
				"uri":      r.RequestURI,
				"method":   r.Method,
				"duration": duration,
				"status":   wrappedWriter.statusCode,
				"size":     wrappedWriter.size,
			}).Info("Handled request")
		})
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}
