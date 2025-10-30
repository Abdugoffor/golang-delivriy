package product_handler

import (
	"log"
	"my-project/helper"
	product_dto "my-project/modul/product/dto"
	product_service "my-project/modul/product/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service product_service.ProductService
}

func NewProductHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) *productHandler {
	handler := &productHandler{
		db:      db,
		log:     log,
		service: product_service.NewProductService(db),
	}
	routes := gorm.Group("/product")
	{
		routes.GET("", handler.All)
		routes.GET("/:id", handler.Show)
		routes.GET("/trash", handler.Trash)
		routes.GET("/trash/:id", handler.ShowTrash)
		routes.POST("", handler.Create)
		routes.PUT("/:id", handler.Update)
		routes.DELETE("/:id", handler.Delete)
		routes.DELETE("/force/:id", handler.ForceDelete)
		routes.PUT("/restore/:id", handler.Restore)
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
		if query.Name != "" {
			tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
		}

		if query.Price != 0 {
			tx = tx.Where("price = ?", query.Price)
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
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	data, err := handler.service.Show(ctx, id)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) Trash(ctx echo.Context) error {
	data, err := handler.service.Trash(ctx)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) ShowTrash(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	data, err := handler.service.ShowTrash(ctx, id)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) Create(ctx echo.Context) error {
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

	return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) Update(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	var req product_dto.Update
	{
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	data, err := handler.service.Update(ctx, id, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *productHandler) Delete(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	err = handler.service.Delete(ctx, id)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success delete data"})
}

func (handler *productHandler) ForceDelete(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	err = handler.service.ForceDelete(ctx, id)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success force delete data"})
}

func (handler *productHandler) Restore(ctx echo.Context) error {
	idParam := ctx.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	data, err := handler.service.Restore(ctx, id)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}
