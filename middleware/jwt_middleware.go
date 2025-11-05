package middleware

import (
	"errors"
	user_model "my-project/modul/user/model"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid token")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token claims"})
			}

			user := &user_model.User{
				ID:    int64(claims["id"].(float64)),
				Email: claims["email"].(string),
			}

			ctx.Set("user", user)
			return next(ctx)
		}
	}
}
