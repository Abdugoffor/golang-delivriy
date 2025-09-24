package helper

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Meta struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}

type PaginatedResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

func Paginate[T any](c echo.Context, db *gorm.DB, model *[]T, count int) (PaginatedResponse[T], error) {
	var res PaginatedResponse[T]

	// 1. Query parameters
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = count
	}

	offset := (page - 1) * limit

	var total int64
	if err := db.Model(model).Count(&total).Error; err != nil {
		return res, err
	}

	// 2. Fetch paginated data
	if err := db.Offset(offset).Limit(limit).Find(model).Error; err != nil {
		return res, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	res = PaginatedResponse[T]{
		Data: *model,
		Meta: Meta{
			Total:       total,
			PerPage:     limit,
			CurrentPage: page,
			LastPage:    lastPage,
		},
	}

	return res, nil
}

func PaginateOnlyTrashed[T any](c echo.Context, db *gorm.DB, model *[]T, count int) (PaginatedResponse[T], error) {
	var res PaginatedResponse[T]

	// 1. Query parameters
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = count
	}

	offset := (page - 1) * limit

	var total int64
	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Model(model).Count(&total).Error; err != nil {
		return res, err
	}

	// 2. Fetch only soft-deleted data
	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Offset(offset).Limit(limit).Find(model).Error; err != nil {
		return res, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	res = PaginatedResponse[T]{
		Data: *model,
		Meta: Meta{
			Total:       total,
			PerPage:     limit,
			CurrentPage: page,
			LastPage:    lastPage,
		},
	}

	return res, nil
}
