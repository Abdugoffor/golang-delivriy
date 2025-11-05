package auth_handler

import (
	"log"
	"my-project/helper"
	"my-project/middleware"
	auth_dto "my-project/modul/auth/dto"
	auth_service "my-project/modul/auth/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service auth_service.AuthService
}

func NewAuthHandler(gorm *echo.Group, db *gorm.DB, log *log.Logger) *authHandler {
	handler := authHandler{
		db:      db,
		log:     log,
		service: auth_service.NewAuthService(db, "secret"),
	}

	routes := gorm.Group("/auth")
	{
		routes.POST("/register", handler.Register)
		routes.POST("/login", handler.Login)
		routes.POST("/refresh", handler.Refresh)
		routes.POST("/logout", handler.Logout)
		routes.GET("/me", handler.Me, middleware.JWTMiddleware("secret"))
	}

	return &handler
}

func (handler *authHandler) Register(ctx echo.Context) error {
	var req auth_dto.RegisterRequest
	{
		if err := ctx.Bind(&req); err != nil {
			return err
		}
	}

	if err := helper.ValidateStruct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	data, err := handler.service.Register(ctx, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusCreated, data)
}

func (handler *authHandler) Login(ctx echo.Context) error {
	var req auth_dto.LoginRequest
	{
		if err := ctx.Bind(&req); err != nil {
			return err
		}
	}

	if err := helper.ValidateStruct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	data, err := handler.service.Login(ctx, req)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *authHandler) Refresh(ctx echo.Context) error {
	data, err := handler.service.Refresh(ctx)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}

func (handler *authHandler) Logout(ctx echo.Context) error {
	if err := handler.service.Logout(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (handler *authHandler) Me(ctx echo.Context) error {
	data, err := handler.service.Me(ctx)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, data)
}
