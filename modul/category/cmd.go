package category_dto

import (
	"log"
	category_handler "my-project/modul/category/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {

	routerGroup := route.Group("/api/v1/admin")
	{
		category_handler.NewCategoryHandler(routerGroup, db, log)
	}

}
