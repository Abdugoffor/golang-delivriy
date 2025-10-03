package user_cmd

import (
	"log"
	user_handler "my-project/modul/user/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {

	routerGroup := route.Group("/api/v1/admin")
	{
		user_handler.NewAuthHandler(routerGroup, db, log)
	}
}
