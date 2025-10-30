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

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productHandler struct {
	db             *gorm.DB
	log            *log.Logger
	service        product_service.ProductService
	historyService history_service.HistoryService
}

func NewProductHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) *productHandler {
	handler := &productHandler{
		db:             db,
		log:            log,
		service:        product_service.NewProductService(db),
		historyService: history_service.NewHistoryService(db),
	}
	routes := gorm.Group("/product")
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

func (handler *productHandler) All(ctx echo.Context) error {
	var query product_dto.Filter
	{
		if err := ctx.Bind(&query); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	filter := func(tx *gorm.DB) *gorm.DB {

		if query.Status == "open" {
			tx = tx.Where("products.deleted_at IS NULL")
		} else if query.Status == "deleted" {
			tx = tx.Where("products.deleted_at IS NOT NULL")
		}

		if query.Name != "" {
			tx = tx.Where("products.name ILIKE ?", "%"+query.Name+"%")
		}

		if query.Price != 0 {
			tx = tx.Where("products.price = ?", query.Price)
		}

		// tx = tx.Group("products.id").
		// 	Order("products.created_at ASC")

		if query.Column != "" && query.Sort != "" {
			groupColumn := fmt.Sprintf("products.%s", query.Column)
			tx = tx.
				Group(groupColumn).
				Order(fmt.Sprintf("products.%s %s", query.Column, query.Sort))
		} else {
			tx = tx.Group("products.id").
				Order("products.created_at DESC")
		}
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

	return helper.View(ctx, "layout.html", "product/index.html", viewData)
}

func (handler *productHandler) Show(ctx echo.Context) error {

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
	return helper.View(ctx, "layout.html", "product/show.html", viewData)
	// return ctx.JSON(http.StatusOK, data)
}
func (handler *productHandler) Edit(ctx echo.Context) error {

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
	return helper.View(ctx, "layout.html", "product/edit.html", viewData)
	// return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) Create(ctx echo.Context) error {

	viewData := map[string]interface{}{}

	return helper.View(ctx, "layout.html", "product/create.html", viewData)
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

	url := fmt.Sprintf("/api/v1/admin/product/%d", data.ID)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (handler *productHandler) Update(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
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

	url := fmt.Sprintf("/api/v1/admin/product/%d", data.ID)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (handler *productHandler) Delete(ctx echo.Context) error {
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

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/product")
}

func (handler *productHandler) ForceDelete(ctx echo.Context) error {
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

	err = handler.service.ForceDelete(ctx, filter)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/product/trash")
}

func (handler *productHandler) Restore(ctx echo.Context) error {
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

	return ctx.Redirect(http.StatusSeeOther, "/api/v1/admin/product")
}
