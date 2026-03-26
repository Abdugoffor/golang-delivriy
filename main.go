package main

import (
	"log"
	"my-project/config"
	"my-project/helper"
	"my-project/middleware"
	category_dto "my-project/modul/category"
	product_cmd "my-project/modul/product"
	"my-project/seeder"

	"github.com/labstack/echo/v4"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	seeder.DBSeed()

	route := echo.New()

	route.Use(middleware.SessionSet())

	product_cmd.Cmd(route, config.DB, log.Default())
	category_dto.Cmd(route, config.DB, log.Default())

	route.Logger.Fatal(route.Start(":" + helper.ENV("HTTP_PORT")))

}
