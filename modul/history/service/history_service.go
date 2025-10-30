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
	var models []history_model.History

	res, err := helper.Paginate(ctx, service.db.Scopes(filter), &models, 5)
	if err != nil {
		return helper.PaginatedResponse[history_dto.Response]{}, err
	}

	var data []history_dto.Response
	for _, model := range models {
		data = append(data, history_dto.ToResponse(model))
	}

	return helper.PaginatedResponse[history_dto.Response]{
		Data: data,
		Meta: res.Meta,
	}, nil
}
