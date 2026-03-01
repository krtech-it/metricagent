package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"slices"
	"strings"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			compressor:     gz,
		}
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer     io.Writer
	compressor *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := w.Header().Get("Content-Type")
	if shouldCompress(contentType) {
		w.Writer = w.compressor
		w.Header().Set("Content-Encoding", "gzip")
	} else {
		w.Writer = w.ResponseWriter
	}
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	contentType := w.Header().Get("Content-Type")
	if shouldCompress(contentType) {
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(code)
}

func shouldCompress(contentType string) bool {
	var contentTypes = [...]string{"application/javascript", "application/json", "text/css", "text/html", "text/plain", "text/xml"}
	if i := strings.Index(contentType, ";"); i != -1 {
		contentType = contentType[:i]
	}
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	return slices.Contains(contentTypes[:], contentType)
}

func DecompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.GetHeader("Content-Encoding")
		if encoding == "" {
			c.Next()
			return
		}

		var reader io.ReadCloser
		switch strings.ToLower(strings.TrimSpace(encoding)) {
		case "gzip":
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip encoding"})
				return
			}
			defer gz.Close()
			reader = gz
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip encoding"})
			return
		}
		c.Request.Body = reader
		c.Next()
	}
}
