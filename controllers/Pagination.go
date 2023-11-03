package controllers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PaginationParams holds the page and limit parameters.
type PaginationParams struct {
	Page  int
	Limit int
}

// ParsePaginationParams parses and returns the page and limit parameters from the query string.
func ParsePaginationParams(c echo.Context) PaginationParams {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page := 1
	limit := 10 // Default limit

	if pageParam != "" {
		pageValue, err := strconv.Atoi(pageParam)
		if err == nil {
			page = pageValue
		}
	}

	if limitParam != "" {
		limitValue, err := strconv.Atoi(limitParam)
		if err == nil {
			limit = limitValue
		}
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetPaginatedData retrieves paginated data from the database based on the provided page and limit.
func GetPaginatedData(c echo.Context, query *gorm.DB, paginationParams PaginationParams, result interface{}) (interface{}, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	if err := query.Limit(paginationParams.Limit).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
