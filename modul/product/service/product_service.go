package product_service

import (
	"my-project/helper"
	product_dto "my-project/modul/product/dto"
	product_model "my-project/modul/product/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductService interface {
	All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[product_dto.Response], error)
	Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (product_dto.Response, error)
	Create(ctx echo.Context, req product_dto.Create) (product_dto.Response, error)
	Update(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB, req product_dto.Update) (product_dto.Response, error)
	Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
}

type productService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return &productService{db: db}
}

func (service *productService) All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[product_dto.Response], error) {
	var models []product_model.Product

	res, err := helper.Paginate(ctx, service.db.Scopes(filter), &models, 10)
	{
		if err != nil {
			return helper.PaginatedResponse[product_dto.Response]{}, err
		}
	}

	var data []product_dto.Response
	{
		for _, model := range models {
			data = append(data, product_dto.ToResponse(model))
		}
	}

	return helper.PaginatedResponse[product_dto.Response]{
		Data: data,
		Meta: res.Meta,
	}, nil
}

func (service *productService) Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (product_dto.Response, error) {
	var model product_model.Product
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return product_dto.Response{}, err
		}
	}

	res := product_dto.ToResponse(model)

	return res, nil
}

func (service *productService) Trash(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[product_dto.Response], error) {
	var models []product_model.Product

	res, err := helper.PaginateOnlyTrashed(ctx, service.db.Scopes(filter), &models, 10)
	{
		if err != nil {
			return helper.PaginatedResponse[product_dto.Response]{}, err
		}
	}

	var data []product_dto.Response
	{
		for _, model := range models {
			data = append(data, product_dto.ToResponse(model))
		}
	}

	return helper.PaginatedResponse[product_dto.Response]{
		Data: data,
		Meta: res.Meta,
	}, nil
}

func (service *productService) ShowTrash(ctx echo.Context, id uint) (product_dto.Response, error) {

	var model product_model.Product
	{
		if err := service.db.Unscoped().Where("id = ?", id).First(&model).Error; err != nil {
			return product_dto.Response{}, err
		}
	}

	res := product_dto.ToResponse(model)

	return res, nil
}

func (service *productService) Create(ctx echo.Context, req product_dto.Create) (product_dto.Response, error) {
	var model product_model.Product
	{
		model.Name = req.Name
		model.Slug = helper.Slug(req.Name)
		model.Price = req.Price

		if err := service.db.Create(&model).Error; err != nil {
			return product_dto.Response{}, err
		}
	}

	res := product_dto.ToResponse(model)

	return res, nil
}

func (service *productService) Update(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB, req product_dto.Update) (product_dto.Response, error) {
	var model product_model.Product
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return product_dto.Response{}, err
		}
	}

	model.Name = req.Name
	model.Price = req.Price

	if err := service.db.Save(&model).Error; err != nil {
		return product_dto.Response{}, err
	}

	res := product_dto.ToResponse(model)

	return res, nil
}

func (service *productService) Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model product_model.Product
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

func (service *productService) ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model product_model.Product
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

func (service *productService) Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error {
	var model product_model.Product
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
