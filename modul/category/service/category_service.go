package category_service

import (
	"my-project/helper"
	category_dto "my-project/modul/category/dto"
	category_model "my-project/modul/category/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryService interface {
	All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[category_dto.Response], error)
	Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (category_dto.Response, error)
	Create(ctx echo.Context, req category_dto.Create) (category_dto.Response, error)
	Update(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB, req category_dto.Update) (category_dto.Response, error)
	Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
}

type categoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{db: db}
}

func (service *categoryService) All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[category_dto.Response], error) {
	return helper.Paginate[category_model.Category, category_dto.Response](ctx, service.db.Scopes(filter), 10)
}

func (service *categoryService) Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (category_dto.Response, error) {
	var model category_model.Category
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return category_dto.Response{}, err
		}
	}

	return category_dto.ToResponse(model), nil
}

func (service *categoryService) Create(ctx echo.Context, req category_dto.Create) (category_dto.Response, error) {
	model := category_model.Category{
		Name:      req.Name,
		Slug:      helper.Slug(req.Name),
		CompanyID: 1,
		IsActive:  req.IsActive,
	}

	if err := service.db.Create(&model).Error; err != nil {
		return category_dto.Response{}, err
	}

	return category_dto.ToResponse(model), nil
}

func (service *categoryService) Update(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB, req category_dto.Update) (category_dto.Response, error) {
	var model category_model.Category
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return category_dto.Response{}, err
		}
	}

	model.Name = req.Name
	model.IsActive = req.IsActive

	if err := service.db.Save(&model).Error; err != nil {
		return category_dto.Response{}, err
	}

	return category_dto.ToResponse(model), nil
}

func (service *categoryService) Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	return service.db.Delete(&model).Error
}

func (service *categoryService) Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Unscoped().Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	return service.db.Model(&model).Unscoped().Update("deleted_at", nil).Error
}

func (service *categoryService) ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Unscoped().Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	return service.db.Unscoped().Delete(&model).Error
}
