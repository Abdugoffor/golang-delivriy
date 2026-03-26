package product_handler

import (
	"fmt"
	"log"
	"my-project/helper"
	history_service "my-project/modul/history/service"
	product_dto "my-project/modul/product/dto"
	product_service "my-project/modul/product/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productHandler struct {
	db             *gorm.DB
	log            *log.Logger
	service        product_service.ProductService
	historyService history_service.HistoryService
}

func NewProductHandler(group *echo.Group, db *gorm.DB, log *log.Logger) *productHandler {
	handler := &productHandler{
		db:             db,
		log:            log,
		service:        product_service.NewProductService(db),
		historyService: history_service.NewHistoryService(db),
	}

	routes := group.Group("/product")
	{
		routes.GET("", handler.All)
		routes.GET("/create", handler.Create)
		routes.GET("/:id", handler.Show)
		routes.GET("/:id/edit", handler.Edit)
		routes.GET("/:id/history", handler.History)
		routes.POST("", handler.Store)
		routes.POST("/:id", handler.Update)
		routes.POST("/:id/delete", handler.Delete)
		routes.POST("/:id/restore", handler.Restore)
		routes.POST("/:id/force", handler.ForceDelete)
	}

	return handler
}

func (handler *productHandler) All(ctx echo.Context) error {
	var query product_dto.Filter
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

		if query.ID != 0 {
			tx = tx.Where("products.id = ?", query.ID)
		}

		if query.Name != "" {
			tx = tx.Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(query.Name)+"%")
		}

		if query.Price != 0 {
			tx = tx.Where("CAST(products.price AS TEXT) LIKE ?", fmt.Sprintf("%%%d%%", query.Price))
		}

		if query.IsActive == "true" {

			tx = tx.Where("products.is_active = ?", true)
		} else if query.IsActive == "false" {

			tx = tx.Where("products.is_active = ?", false)
		}

		orderClause := "products.created_at ASC"
		if query.Column != "" && query.Sort != "" {
			orderClause = fmt.Sprintf("products.%s %s", query.Column, query.Sort)
		}

		return tx.Group("products.id").Order(orderClause)
	}

	data, err := handler.service.All(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return helper.View(ctx, "layout.html", "product/index.html", map[string]any{
		"Title":  "Товары",
		"models": data.Data,
		"Meta":   data.Meta,
		"Filter": query,
	})
}

func (handler *productHandler) History(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	modelFilter := func(tx *gorm.DB) *gorm.DB {
		return tx.Unscoped().Where("id = ?", id)
	}

	model, err := handler.service.Show(ctx, modelFilter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	historyFilter := func(tx *gorm.DB) *gorm.DB {
		return tx.Where("table_name = ? AND row_id = ?", "products", id)
	}

	data, err := handler.historyService.All(ctx, historyFilter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return helper.View(ctx, "layout.html", "product/history.html", map[string]any{
		"Title":  "История: " + model.Name,
		"model":  model,
		"models": data.Data,
		"Meta":   data.Meta,
	})
}

func (handler *productHandler) Show(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	data, err := handler.service.Show(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return helper.View(ctx, "layout.html", "product/show.html", map[string]any{
		"Title": data.Name,
		"model": data,
	})
}

func (handler *productHandler) Create(ctx echo.Context) error {
	return helper.View(ctx, "layout.html", "product/create.html", map[string]any{
		"Title": "Товары",
	})
}

func (handler *productHandler) Store(ctx echo.Context) error {
	var req product_dto.Create
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

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/api/v1/admin/product/%d", data.ID))
}

func (handler *productHandler) Edit(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	data, err := handler.service.Show(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return helper.View(ctx, "layout.html", "product/edit.html", map[string]any{
		"Title": "Edit: " + data.Name,
		"model": data,
	})
}

func (handler *productHandler) Update(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	var req product_dto.Update
	{
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	data, err := handler.service.Update(ctx, filter, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/api/v1/admin/product/%d", data.ID))
}

func (handler *productHandler) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	if err := handler.service.Delete(ctx, filter); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	back := ctx.FormValue("_back")
	if !strings.HasPrefix(back, "/api/v1/admin/") {
		back = "/api/v1/admin/product"
	}
	return ctx.Redirect(http.StatusSeeOther, back)
}

func (handler *productHandler) Restore(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	if err := handler.service.Restore(ctx, filter); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	back := ctx.FormValue("_back")
	if !strings.HasPrefix(back, "/api/v1/admin/") {
		back = "/api/v1/admin/product"
	}
	return ctx.Redirect(http.StatusSeeOther, back)
}

func (handler *productHandler) ForceDelete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {
		tx.Where("id = ?", id)
		return tx
	}

	if err := handler.service.ForceDelete(ctx, filter); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	back := ctx.FormValue("_back")
	if !strings.HasPrefix(back, "/api/v1/admin/") {
		back = "/api/v1/admin/product"
	}
	return ctx.Redirect(http.StatusSeeOther, back)
}
