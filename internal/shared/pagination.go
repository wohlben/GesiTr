package shared

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultLimit = 50
	MaxLimit     = 200
)

type PaginationParams struct {
	Limit  int
	Offset int
}

func ParsePagination(c *gin.Context) PaginationParams {
	p := PaginationParams{Limit: DefaultLimit}

	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			p.Limit = n
		}
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}

	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			p.Offset = n
		}
	}

	return p
}

func ApplyPagination(db *gorm.DB, p PaginationParams) *gorm.DB {
	return db.Limit(p.Limit).Offset(p.Offset)
}
