package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func GzipDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			contentType := r.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
				gzReader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, "Invalid gzip data", http.StatusBadRequest)
					return
				}
				defer gzReader.Close()

				r.Body = io.NopCloser(gzReader)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GzipCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzWriter := gzip.NewWriter(w)
			defer gzWriter.Close()

			gzipWriter := &gzipResponseWriter{ResponseWriter: w, writer: gzWriter}
			next.ServeHTTP(gzipWriter, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
