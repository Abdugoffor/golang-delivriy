package auth_cmd

import (
	"log"
	auth_handler "my-project/modul/auth/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {

	routerGroup := route.Group("/api/v1/admin")
	{
		auth_handler.NewAuthHandler(routerGroup, db, log)
	}
}
