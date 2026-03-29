package auth

import (
	"log"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestLogLevel returns the configured REQUEST_LOG level: "trace", "info", or "".
func RequestLogLevel() string {
	return strings.ToLower(os.Getenv("REQUEST_LOG"))
}

// RequestTrace returns a Gin middleware that dumps the full incoming request
// when REQUEST_LOG=trace. Place before auth middleware to see headers as they
// arrive from the reverse proxy / client.
func RequestTrace() gin.HandlerFunc {
	level := RequestLogLevel()

	return func(c *gin.Context) {
		if level == "trace" {
			dump, err := httputil.DumpRequest(c.Request, true)
			if err != nil {
				log.Printf("[TRACE] failed to dump request: %v", err)
			} else {
				log.Printf("[TRACE] %s %s\n%s", c.Request.Method, c.Request.URL.Path, dump)
			}
		}
		c.Next()
	}
}

// APIOnlyLogger returns a Gin logger that only logs /api requests,
// filtering out static resource noise. Used when REQUEST_LOG=info or trace.
func APIOnlyLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: nil,
		Skip: func(c *gin.Context) bool {
			return !strings.HasPrefix(c.Request.URL.Path, "/api")
		},
	})
}
