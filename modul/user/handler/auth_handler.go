package user_handler

import (
	"log"
	middleware "my-project/middlewares"
	user_dto "my-project/modul/user/dto"
	auth_service "my-project/modul/user/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db      *gorm.DB
	log     *log.Logger
	service auth_service.AuthService
}

func NewAuthHandler(g *echo.Group, db *gorm.DB, log *log.Logger) *AuthHandler {
	h := &AuthHandler{
		db:      db,
		log:     log,
		service: auth_service.NewAuthService(db),
	}

	r := g.Group("/auth")
	{
		r.POST("/register", h.Register)
		r.GET("/verify", h.VerifyEmail)
		r.POST("/login", h.Login)
		r.POST("/logout", h.Logout)
		r.GET("/me", h.Profile, middleware.AuthMiddleware())
		r.PUT("/me", h.UpdateProfile, middleware.AuthMiddleware())
		r.POST("/password-reset/request", h.RequestPasswordReset)
		r.POST("/password-reset", h.ResetPassword)
	}
	return h
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req user_dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error 1": err.Error()})
	}
	// validation optional - you can use validator lib
	user, err := h.service.Register(&auth_service.RegisterDto{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error 2": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "registered", "user_id": user.ID})
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "token required"})
	}
	if err := h.service.VerifyEmail(token); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "email verified"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req user_dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user_dto.LoginResponse{Token: token})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// For stateless JWT, logout is typically performed client-side by deleting token.
	// If you want to invalidate tokens server-side, implement token blacklist table.
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}

func (h *AuthHandler) Profile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.service.FindById(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"is_verified": user.IsVerified,
	})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req user_dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	user, err := h.service.UpdateProfile(userID, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "updated", "user": echo.Map{
		"id": user.ID, "name": user.Name, "email": user.Email,
	}})
}

func (h *AuthHandler) RequestPasswordReset(c echo.Context) error {
	var req user_dto.RequestPasswordReset
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	_ = h.service.RequestPasswordReset(req.Email)
	// always return ok to avoid revealing whether email exists
	return c.JSON(http.StatusOK, echo.Map{"message": "if email exists, reset link sent"})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req user_dto.ResetPassword
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "password reset success"})
}
