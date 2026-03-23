package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/config"
	"io"
	"net/http"
)

func CheckHashMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientHash := c.Request.Header.Get("HashSHA256")
		if cfg.HashKey == nil || clientHash == "" {
			c.Next()
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		h := hmac.New(sha256.New, []byte(*cfg.HashKey))
		h.Write(body)
		serverHash := hex.EncodeToString(h.Sum(nil))
		if serverHash != clientHash {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid hash"})
			return
		}
		c.Next()
	}
}

func ResponseHashMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.HashKey == nil {
			c.Next()
			return
		}
		c.Writer = &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			hashKey:        *cfg.HashKey,
		}
		c.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	hashKey string
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	h := hmac.New(sha256.New, []byte(w.hashKey))
	h.Write(w.body.Bytes())
	serverHash := hex.EncodeToString(h.Sum(nil))
	w.Header().Set("HashSHA256", serverHash)
	return w.ResponseWriter.Write(b)
}
