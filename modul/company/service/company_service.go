package company_service

import (
	"my-project/helper"
	company_dto "my-project/modul/company/dto"
	company_model "my-project/modul/company/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CompanyService interface {
	All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[company_dto.Response], error)
	Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (company_dto.Response, error)
	// Create(ctx echo.Context, req category_dto.Create) (category_dto.Response, error)
	// Update(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB, req category_dto.Update) (category_dto.Response, error)
	// Delete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	// Restore(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
	// ForceDelete(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) error
}

type companyService struct {
	db *gorm.DB
}

func NewCompanyService(db *gorm.DB) CompanyService {
	return &companyService{db: db}
}
func (service *companyService) All(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (helper.PaginatedResponse[company_dto.Response], error) {
	var models []company_model.Company

	res, err := helper.Paginate(ctx, service.db.Scopes(filter), &models, 10)
	{
		if err != nil {
			return helper.PaginatedResponse[company_dto.Response]{}, err
		}
	}
	var data []company_dto.Response
	{
		for _, model := range models {
			data = append(data, company_dto.ToResponse(model))
		}
	}

	return helper.PaginatedResponse[company_dto.Response]{
		Data: data,
		Meta: res.Meta,
	}, nil
}

func (service *companyService) Show(ctx echo.Context, filter func(tx *gorm.DB) *gorm.DB) (company_dto.Response, error) {
	var model company_model.Company
	{
		if err := service.db.Scopes(filter).First(&model).Error; err != nil {
			return company_dto.Response{}, err
		}
	}

	res := company_dto.ToResponse(model)

	return res, nil
}
