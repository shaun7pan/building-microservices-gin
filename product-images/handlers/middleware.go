package handlers

import (
	"compress/gzip"
	"fmt"

	"github.com/gin-gonic/gin"
)

type wrappedGinWriter struct {
	gin.ResponseWriter
	gw *gzip.Writer
}

// Custom middleware
func (f *Files) GzipMiddleware(c *gin.Context) {
	f.log.Info("gzip middleware")
	if c.Request.Header.Get("Accept-Encoding") == "gzip" {
		//create a gziped response
		wgw := NewWrappedGinWriter(c.Writer)
		wgw.Header().Set("Content-Encoding", "gzip")
		c.Writer = wgw
		defer func() {
			wgw.Flush()
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()
	}
	c.Next()
}

func NewWrappedGinWriter(rw gin.ResponseWriter) *wrappedGinWriter {
	return &wrappedGinWriter{rw, gzip.NewWriter(rw)}
}

func (wr *wrappedGinWriter) Write(d []byte) (int, error) {
	wr.Header().Del("Content-Length")
	return wr.gw.Write(d)
}

func (wr *wrappedGinWriter) Flush() {

	wr.gw.Flush()
	wr.gw.Close()
}
