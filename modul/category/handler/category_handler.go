package category_handler

import (
	"fmt"
	"log"
	"my-project/helper"
	category_dto "my-project/modul/category/dto"
	category_service "my-project/modul/category/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type categoryHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service category_service.CategoryService
}

func NewCategoryHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) categoryHandler {
	handler := categoryHandler{
		db:      db,
		log:     log,
		service: category_service.NewCategoryService(db),
	}

	routes := gorm.Group("/category")
	{
		routes.GET("", handler.All)
		routes.GET("/:id", handler.Show)
		routes.GET("/create", handler.Create)
		routes.POST("", handler.Store)
		routes.GET("/:id/edit", handler.Edit)
		routes.POST("/:id", handler.Update)
		routes.POST("/:id/delete", handler.Delete)
		routes.POST("/:id/restore", handler.Restore)
		routes.POST("/:id/force", handler.ForceDelete)
	}

	return handler
}

func (handler *categoryHandler) All(ctx echo.Context) error {
	var query category_dto.Filter
	{
		if err := ctx.Bind(&query); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		switch query.Status {
		case "open":
			tx = tx.Where("deleted_at IS NULL")
		case "deleted":
			tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
		default:
			tx = tx.Unscoped()
		}

		if query.Name != "" {
			name := "%" + strings.ToLower(query.Name) + "%"
			tx = tx.Where("LOWER(categories.name) LIKE ?", name)
		}

		orderClause := "categories.created_at ASC"
		if query.Column != "" && query.Sort != "" {
			orderClause = fmt.Sprintf("categories.%s %s", query.Column, query.Sort)
		}

		tx = tx.Group("categories.id").Order(orderClause)

		return tx
	}

	data, err := handler.service.All(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	viewData := map[string]interface{}{
		"models": data.Data,
		"Meta":   data.Meta,
		"Filter": query,
	}

	return helper.View(ctx, "layout.html", "category/index.html", viewData)
}

func (handler *categoryHandler) Show(ctx echo.Context) error {

	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseInt(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	data, err := handler.service.Show(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}
	viewData := map[string]interface{}{
		"model": data,
	}
	return helper.View(ctx, "layout.html", "category/show.html", viewData)
}

func (handler *categoryHandler) Create(ctx echo.Context) error {

	viewData := map[string]interface{}{}

	return helper.View(ctx, "layout.html", "category/create.html", viewData)
}

func (handler *categoryHandler) Store(ctx echo.Context) error {
	var req category_dto.Create
	{
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	data, err := handler.service.Create(ctx, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	url := fmt.Sprintf("/api/v1/admin/category/%d", data.ID)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (handler *categoryHandler) Edit(ctx echo.Context) error {

	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseInt(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	data, err := handler.service.Show(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}
	viewData := map[string]interface{}{
		"model": data,
	}
	return helper.View(ctx, "layout.html", "category/edit.html", viewData)
}

func (handler *categoryHandler) Update(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	var req category_dto.Update
	{
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	data, err := handler.service.Update(ctx, filter, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	url := fmt.Sprintf("/api/v1/admin/category/%d", data.ID)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (handler *categoryHandler) Delete(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	err = handler.service.Delete(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/category")
}

func (handler *categoryHandler) Restore(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	if err := handler.service.Restore(ctx, filter); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/category")
}

func (handler *categoryHandler) ForceDelete(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if parsedID > 0 {
			tx = tx.Where("id = ?", parsedID)
		}
		return tx
	}

	if err := handler.service.ForceDelete(ctx, filter); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/category")
}
