package history

import (
	"github.com/labstack/echo/v4"
)

// RequestContext middleware
// Har bir request’da context’ga user_id, ip va api path yoziladi
func RequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Masalan JWT’dan user_id olish mumkin, hozircha 1 deb qoldiramiz
		c.Set("user_id", int64(1))

		ip := c.RealIP()
		c.Set("ip", ip)

		api := c.Request().URL.Path
		c.Set("api", api)

		return next(c)
	}
}
