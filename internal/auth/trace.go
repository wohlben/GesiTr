package auth

import (
	"log"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestTrace returns a Gin middleware that dumps the full incoming request
// when REQUEST_LOG=trace. Place before auth middleware to see headers as they
// arrive from the reverse proxy / client.
func RequestTrace() gin.HandlerFunc {
	enabled := strings.EqualFold(os.Getenv("REQUEST_LOG"), "trace")

	return func(c *gin.Context) {
		if enabled {
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
