package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

type paginationParams struct {
	Limit  int
	Offset int
}

func parsePagination(c *gin.Context) paginationParams {
	p := paginationParams{Limit: defaultLimit}

	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			p.Limit = n
		}
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			p.Offset = n
		}
	}

	return p
}

func applyPagination(db *gorm.DB, p paginationParams) *gorm.DB {
	return db.Limit(p.Limit).Offset(p.Offset)
}
