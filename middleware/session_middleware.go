package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sess, _ := session.Get("session", ctx)

		token, ok := sess.Values["token"].(string)
		{
			if !ok || token == "" {
				return ctx.Redirect(http.StatusFound, "/api/v1/admin/auth/login")
			}
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return ctx.Redirect(http.StatusFound, "/api/v1/admin/auth/login")
		}

		ctx.Set("user_id", claims["id"])
		ctx.Set("user_email", claims["email"])
		return next(ctx)
	}
}
