package middleware

import (
	"my-project/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization"})
			}

			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization header"})
			}

			token := strings.TrimPrefix(auth, "Bearer ")

			claims, err := helper.JwtParse(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			// user_id claimni olish
			switch v := claims["user_id"].(type) {
			case float64:
				c.Set("user_id", uint(uint64(v)))
			case string:
				uid, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token claims"})
				}
				c.Set("user_id", uint(uid))
			default:
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token claims"})
			}

			return next(c)
		}
	}
}
