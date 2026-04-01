package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationQuery struct {
	Page     int
	PageSize int
}

func ParsePagination(c *gin.Context) (PaginationQuery, error) {
	result := PaginationQuery{}
	var err error

	if raw := c.Query("page"); raw != "" {
		result.Page, err = strconv.Atoi(raw)
		if err != nil {
			return PaginationQuery{}, err
		}
	}
	if raw := c.Query("page_size"); raw != "" {
		result.PageSize, err = strconv.Atoi(raw)
		if err != nil {
			return PaginationQuery{}, err
		}
	}

	return result, nil
}
