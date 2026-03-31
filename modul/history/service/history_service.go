package history_service

import (
	"my-project/helper"
	history_dto "my-project/modul/history/dto"
	history_model "my-project/modul/history/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HistoryService interface {
	All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[history_dto.Response], error)
}

type historyService struct {
	db *gorm.DB
}

func NewHistoryService(db *gorm.DB) HistoryService {
	return &historyService{db: db}
}

func (service *historyService) All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[history_dto.Response], error) {
	return helper.Paginate[history_model.History, history_dto.Response](ctx, service.db.Scopes(filter).Order("id DESC"), 5)
}
