package app_cmd

import (
	"log"
	app_handler "my-project/modul/app/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {
	routerGroup := route.Group("")
	{
		app_handler.NewAppHandler(routerGroup, db, log)
		app_handler.NewAppCateHandler(routerGroup, db, log)
	}
}
