package category_handler

import (
	"log"
	category_dto "my-project/modul/category/dto"
	category_service "my-project/modul/category/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type categoryHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service category_service.CategoryService
}

func NewCategoryHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) *categoryHandler {
	handler := &categoryHandler{
		db:      db,
		log:     log,
		service: category_service.NewCategoryService(db),
	}
	routes := gorm.Group("/category")
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

func (handler *categoryHandler) All(ctx echo.Context) error {

	data, err := handler.service.All(ctx)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *categoryHandler) Show(ctx echo.Context) error {

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

func (handler *categoryHandler) Trash(ctx echo.Context) error {

	data, err := handler.service.Trash(ctx)
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *categoryHandler) ShowTrash(ctx echo.Context) error {

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

func (handler *categoryHandler) Create(ctx echo.Context) error {

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

	return ctx.JSON(http.StatusOK, data)
}

func (handler *categoryHandler) Update(ctx echo.Context) error {

	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}
	}

	id := uint(parsedID)

	var req category_dto.Update
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

func (handler *categoryHandler) Delete(ctx echo.Context) error {

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

func (handler *categoryHandler) ForceDelete(ctx echo.Context) error {

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

func (handler *categoryHandler) Restore(ctx echo.Context) error {

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
