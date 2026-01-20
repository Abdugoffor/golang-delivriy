package app_service

import (
	"context"
	app_dto "my-project/modul/app/dto"
	app_model "my-project/modul/app/model"

	"github.com/Abdugoffor/echo-crud-pg/pg"
	"github.com/Abdugoffor/echo-crud-pg/request"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppService interface {
	All(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*app_dto.AppPage, error)
	Show(ctx context.Context, filter pg.Filter) (*app_dto.AppResponse, error)
	Page(ctx context.Context, filter pg.Filter) ([]app_dto.AppResponse, error)
	CreateCate(ctx echo.Context, req app_dto.CreateCate) (int64, error)
	Create(ctx echo.Context, req app_dto.Create) (int64, error)
}

type appService struct {
	db *gorm.DB
}

func NewAppService(db *gorm.DB) AppService {
	return &appService{db: db}
}

func (service *appService) All(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*app_dto.AppPage, error) {
	return pg.PageWithScan[app_model.App, app_dto.AppResponse](service.db, paginate, filter)
}

func (service *appService) Page(ctx context.Context, filter pg.Filter) ([]app_dto.AppResponse, error) {
	return pg.FindWithScan[app_model.App, app_dto.AppResponse](service.db, filter)
}
func (service *appService) Show(ctx context.Context, filter pg.Filter) (*app_dto.AppResponse, error) {
	return pg.FindOneWithScan[app_model.App, app_dto.AppResponse](service.db, filter)
}

func (service *appService) Create(ctx echo.Context, req app_dto.Create) (int64, error) {
	model := app_model.App{
		UserID:        2,
		AppCategoryID: req.AppCategoryID,
	}

	if err := service.db.Create(&model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (service *appService) CreateCate(ctx echo.Context, req app_dto.CreateCate) (int64, error) {

	for _, f := range req.AppFormField {
		model := app_model.AppValue{
			AppID:     req.AppID,
			AppPageID: f.AppPageID,
			AppFormID: f.AppFormID,
			Value:     f.Value,
			IsActive:  true,
		}

		if err := service.db.Create(&model).Error; err != nil {
			return 0, err
		}
	}

	return req.AppID, nil
}
