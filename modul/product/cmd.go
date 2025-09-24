package product_cmd

import (
	"log"
	product_handler "my-project/modul/product/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {

	routerGroup := route.Group("/api/v1/admin")
	{
		product_handler.NewProductHandler(routerGroup, db, log)
	}
}
