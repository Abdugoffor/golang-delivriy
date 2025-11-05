package company_cmd

import (
	"log"
	company_handler "my-project/modul/company/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(route *echo.Echo, db *gorm.DB, log *log.Logger) {

	routerGroup := route.Group("/api/v1/admin")
	{
		company_handler.NewCompanyHandler(routerGroup, db, log)
	}

}
