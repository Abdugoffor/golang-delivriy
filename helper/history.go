package helper

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterHistoryCallbacks(db *gorm.DB, e *echo.Echo) {
	// Create
	db.Callback().Create().Before("gorm:create").Register("history:set_local", func(tx *gorm.DB) {
		setLocalVars(tx)
	})
	// Update
	db.Callback().Update().Before("gorm:update").Register("history:set_local", func(tx *gorm.DB) {
		setLocalVars(tx)
	})
	// Delete
	db.Callback().Delete().Before("gorm:delete").Register("history:set_local", func(tx *gorm.DB) {
		setLocalVars(tx)
	})
}

func setLocalVars(tx *gorm.DB) {
	// Echo contextni olish
	c, ok := tx.Statement.Context.Value("echo_context").(echo.Context)
	if !ok {
		return
	}

	// Contextdan olingan ma’lumotlar
	userID, _ := c.Get("history_user_id").(string)
	ip, _ := c.Get("history_ip").(string)
	path, _ := c.Get("history_path").(string)
	method, _ := c.Get("history_method").(string)

	// 🔑 Postgres session variables (SET LOCAL)
	tx.Exec(fmt.Sprintf("SET LOCAL app.current_user_id = '%s'", userID))
	tx.Exec(fmt.Sprintf("SET LOCAL app.current_request_ip = '%s'", ip))
	tx.Exec(fmt.Sprintf("SET LOCAL app.current_request_path = '%s'", path))
	tx.Exec(fmt.Sprintf("SET LOCAL app.current_request_method = '%s'", method))
}
