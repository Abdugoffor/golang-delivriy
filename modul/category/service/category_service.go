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
	var models []category_model.Category

	res, err := helper.Paginate(ctx, service.db.Scopes(filter), &models, 10)
	{
		if err != nil {
			return helper.PaginatedResponse[category_dto.Response]{}, err
		}
	}

	var data []category_dto.Response
	{
		for _, model := range models {
			data = append(data, category_dto.ToResponse(model))
		}
	}

	return helper.PaginatedResponse[category_dto.Response]{
		Data: data,
		Meta: res.Meta,
	}, nil

}

func (service *categoryService) Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (category_dto.Response, error) {
	var model category_model.Category
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return category_dto.Response{}, err
		}
	}

	res := category_dto.ToResponse(model)

	return res, nil
}

func (service *categoryService) Create(ctx echo.Context, req category_dto.Create) (category_dto.Response, error) {
	var model category_model.Category
	{
		model.Name = req.Name
		model.Slug = helper.Slug(req.Name)
		model.CompanyID = 1
		model.IsActive = req.IsActive

		if err := service.db.Create(&model).Error; err != nil {
			return category_dto.Response{}, err
		}
	}

	res := category_dto.ToResponse(model)

	return res, nil
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

	res := category_dto.ToResponse(model)

	return res, nil
}

func (service *categoryService) Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	if err := service.db.Delete(&model).Error; err != nil {
		return err
	}

	return nil
}

func (service *categoryService) Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Unscoped().Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	if err := service.db.Model(&model).Unscoped().Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

func (service *categoryService) ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model category_model.Category
	{
		if err := service.db.Unscoped().Scopes(filter).First(&model).Error; err != nil {
			return err
		}
	}

	if err := service.db.Unscoped().Delete(&model).Error; err != nil {
		return err
	}

	return nil
}
