package humaconfig

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// NewAPI creates a huma API that shares the given Gin router group.
// Routes registered on this API go through the group's middleware (auth, profile).
// Paths in huma operations should be relative to the group (e.g., "/exercises"
// not "/api/exercises"). The OpenAPI spec uses a server URL of "/api".
func NewAPI(r *gin.Engine, group *gin.RouterGroup) huma.API {
	config := huma.DefaultConfig("GesiTr API", "1.0.0")
	config.Servers = []*huma.Server{{URL: "/api"}}
	return humagin.NewWithGroup(r, group, config)
}
