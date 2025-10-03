package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

// HistoryMiddleware - user va request ma'lumotlarini Postgres kontekstiga yozib qo‘yadi
func HistoryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// User ID
			userID := "0"
			if uid := c.Request().Header.Get("X-User-ID"); uid != "" {
				userID = uid
			}

			// Request info
			ip := c.RealIP()
			path := c.Path()
			method := c.Request().Method

			// Contextga yozib qo‘yamiz (echo context uchun)
			c.Set("history_user_id", userID)
			c.Set("history_ip", ip)
			c.Set("history_path", path)
			c.Set("history_method", method)

			// 🔑 GORM uchun ham contextga echo.Context qo‘shib yuboramiz
			reqCtx := context.WithValue(c.Request().Context(), "echo_context", c)
			c.SetRequest(c.Request().WithContext(reqCtx))

			return next(c)
		}
	}
}
