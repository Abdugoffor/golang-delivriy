package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// SessionSet - session middleware'ni qaytaradi
func SessionSet() echo.MiddlewareFunc {
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,     // HTTPS bo'lsa true qilinadi
		MaxAge:   86400 * 7, // 7 kunlik session
	}

	return session.Middleware(store)
}
