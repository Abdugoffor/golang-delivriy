package app_service

import (
	"context"
	app_dto "my-project/modul/app/dto"
	app_model "my-project/modul/app/model"

	"github.com/Abdugoffor/echo-crud-pg/pg"
	"github.com/Abdugoffor/echo-crud-pg/request"
	"gorm.io/gorm"
)

type AppService interface {
	All(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*app_dto.AppPage, error)
	Show(ctx context.Context, filter pg.Filter) (*app_dto.AppResponse, error)
	Page(ctx context.Context, filter pg.Filter) ([]app_dto.AppResponse, error)
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
