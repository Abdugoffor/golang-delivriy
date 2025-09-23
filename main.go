package main

import (
	"my-project/config"
	"my-project/helper"
	"my-project/seeder"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	seeder.DBSeed()

	route := echo.New()

	route.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Salom, Echo ishlayapti 🚀")
	})

	route.Logger.Fatal(route.Start(":" + helper.ENV("HTTP_PORT")))

}
