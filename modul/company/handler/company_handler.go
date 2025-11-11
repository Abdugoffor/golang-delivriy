package company_handler

import (
	"fmt"
	"log"
	"my-project/helper"
	"my-project/middleware"
	company_dto "my-project/modul/company/dto"
	company_service "my-project/modul/company/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type companyHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service company_service.CompanyService
}

func NewCompanyHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) companyHandler {
	handler := companyHandler{
		db:      db,
		log:     log,
		service: company_service.NewCompanyService(db),
	}

	routes := gorm.Group("/company", middleware.SessionAuthMiddleware)
	{
		routes.GET("", handler.All)
		routes.GET("/:id", handler.Show)
		// routes.GET("/create", handler.Create)
		// routes.POST("", handler.Store)
		// routes.GET("/:id/edit", handler.Edit)
		// routes.POST("/:id", handler.Update)
		// routes.POST("/:id/delete", handler.Delete)
		// routes.POST("/:id/restore", handler.Restore)
		// routes.POST("/:id/force", handler.ForceDelete)
	}

	return handler
}

func (handler *companyHandler) All(ctx echo.Context) error {
	var query company_dto.Filter
	{
		if err := ctx.Bind(&query); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	user := helper.AuthUser(ctx)

	if user == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login")
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		switch query.Status {
		case "open":
			tx = tx.Where("company.deleted_at IS NULL")
		case "deleted":
			tx = tx.Unscoped().Where("company.deleted_at IS NOT NULL")
		default:
			tx = tx.Unscoped()
		}

		if query.Name != "" {
			name := "%" + strings.ToLower(query.Name) + "%"
			tx = tx.Where("LOWER(company.name) LIKE ?", name)
		}

		tx = tx.Preload("Categories").
			Preload("CompanyUserRoles").
			Preload("CompanyUserRoles.User").
			Preload("CompanyUserRoles.Role")

		orderClause := "company.created_at ASC"
		{
			if query.Column != "" && query.Sort != "" {
				orderClause = fmt.Sprintf("company.%s %s", query.Column, query.Sort)
			}
		}

		
		tx = tx.Joins("JOIN company_user_roles ON company_user_roles.company_id = company.id").
			Where("company_user_roles.user_id = ?", user["id"]).
			Group("company.id").
			Order(orderClause)

		return tx
	}

	data, err := handler.service.All(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)

	// viewData := map[string]interface{}{
	// 	"models": data.Data,
	// 	"Meta":   data.Meta,
	// 	"Filter": query,
	// }

	// return helper.View(ctx, "layout.html", "company/index.html", viewData)
}

func (handler *companyHandler) Show(ctx echo.Context) error {
	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseInt(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		tx = tx.Preload("Categories").
			Preload("CompanyUserRoles").
			Preload("CompanyUserRoles.User").
			Preload("CompanyUserRoles.Role")

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}

		tx = tx.Group("company.id")

		return tx
	}

	data, err := handler.service.Show(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)

	// viewData := map[string]interface{}{
	// 	"models": data.Data,
	// 	"Meta":   data.Meta,
	// 	"Filter": query,
	// }

	// return helper.View(ctx, "layout.html", "company/index.html", viewData)
}
