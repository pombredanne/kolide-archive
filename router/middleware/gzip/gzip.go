package gzip

import (
	"compress/gzip"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// BestCompression gzip mode
	BestCompression = gzip.BestCompression
	// BestSpeed gzip mode
	BestSpeed = gzip.BestSpeed
	// DefaultCompression gzip mode
	DefaultCompression = gzip.DefaultCompression
	// NoCompression gzip mode
	NoCompression = gzip.NoCompression
)

// Gzip request data
func Gzip(level int) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !shouldCompress(c.Request) {
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, level)

		if err != nil {
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		c.Writer = &gzipWriter{c.Writer, gz}

		defer func() {
			c.Header("Content-Length", "")
			gz.Close()
		}()

		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	if strings.Contains(req.Header.Get("Connection"), "Upgrade") {
		return false
	}

	extension := filepath.Ext(req.URL.Path)

	// fast path
	if len(extension) < 4 {
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
