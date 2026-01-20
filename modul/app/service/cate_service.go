package app_service

import (
	"context"
	app_dto "my-project/modul/app/dto"
	app_model "my-project/modul/app/model"

	"github.com/Abdugoffor/echo-crud-pg/pg"
	"github.com/Abdugoffor/echo-crud-pg/request"
	"gorm.io/gorm"
)

type AppCate interface {
	All(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*app_dto.AppCatePage, error)
	Show(ctx context.Context, filter pg.Filter) (*app_dto.AppCateResponse, error)
}

type appCate struct {
	db *gorm.DB
}

func NewAppCate(db *gorm.DB) AppCate {
	return &appCate{db: db}
}

func (service *appCate) All(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*app_dto.AppCatePage, error) {
	return pg.PageWithScan[app_model.AppCategory, app_dto.AppCateResponse](service.db, paginate, filter)
}

func (service *appCate) Show(ctx context.Context, filter pg.Filter) (*app_dto.AppCateResponse, error) {
	return pg.FindOneWithScan[app_model.AppCategory, app_dto.AppCateResponse](service.db, filter)
}
