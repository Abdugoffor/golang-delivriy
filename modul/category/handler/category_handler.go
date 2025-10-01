package category_handler

import (
	"fmt"
	"log"
	category_dto "my-project/modul/category/dto"
	category_service "my-project/modul/category/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nguyenthenguyen/docx"
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
	// 1. URL param olish
	idParam := ctx.Param("id")

	parsedID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	id := uint(parsedID)

	// 2. Serviceden ma'lumot olish
	data, err := handler.service.Show(ctx, id)
	if err != nil {
		handler.log.Println("service error:", err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	// 3. Executable pathni olish (aniq project rootni topish uchun)
	exe, _ := os.Executable()
	exePath := filepath.Dir(exe)

	// 4. Template fayl yo‘lini yasash
	templatePath := filepath.Join(exePath, "templates", "template2.docx")
	handler.log.Println("Template path:", templatePath)

	// 5. DOCX faylni ochish
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		handler.log.Println("template open error:", err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "template file not found"})
	}
	defer r.Close()

	// 6. Editable doc olish
	docx1 := r.Editable()

	// 7. Placeholderlarni almashtirish
	docx1.Replace("{{name}}", data.Name, -1)
	docx1.Replace("{{surname}}", data.Slug, -1)
	docx1.Replace("{{date}}", time.Now().Format("2006-01-02"), -1)

	// 8. Output faylni vaqt asosida nomlash
	outputFile := filepath.Join(exePath, fmt.Sprintf("generated_%d.docx", time.Now().UnixNano()))

	// 9. DOCX faylni yozish
	if err := docx1.WriteToFile(outputFile); err != nil {
		handler.log.Println("write error:", err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to write file"})
	}

	// 10. Foydalanuvchiga faylni yuklab berish
	return ctx.Attachment(outputFile, "category.docx")
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
