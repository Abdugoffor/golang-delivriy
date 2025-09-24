package main

import (
	"log"
	"my-project/config"
	"my-project/helper"
	category_cmd "my-project/modul/category"
	"my-project/seeder"

	"github.com/labstack/echo/v4"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	seeder.DBSeed()

	route := echo.New()

	category_cmd.Cmd(route, config.DB, log.Default())

	route.Logger.Fatal(route.Start(":" + helper.ENV("HTTP_PORT")))

}
