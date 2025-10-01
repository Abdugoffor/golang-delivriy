package history

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RegisterHooks – barcha Create/Update/Delete/Restore eventlariga hook ulash
func RegisterHooks(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("history:create", func(tx *gorm.DB) {
		storeHistory(tx, "create", nil, tx.Statement.Dest)
	})

	db.Callback().Update().After("gorm:update").Register("history:update", func(tx *gorm.DB) {
		if tx.Statement.Changed() {
			storeHistory(tx, "update", tx.Statement.ReflectValue.Interface(), tx.Statement.Dest)
		}
	})

	db.Callback().Delete().After("gorm:delete").Register("history:delete", func(tx *gorm.DB) {
		storeHistory(tx, "delete", tx.Statement.Dest, nil)
	})

	// restore uchun (gorm soft delete)
	db.Callback().Update().After("gorm:restore").Register("history:restore", func(tx *gorm.DB) {
		storeHistory(tx, "restore", nil, tx.Statement.Dest)
	})
}

// storeHistory – umumiy yozuvchi
func storeHistory(tx *gorm.DB, action string, oldModel interface{}, newModel interface{}) {
	table := tx.Statement.Table
	rowID := fmt.Sprintf("%v", getPrimaryKey(newModel))

	var oldJSON, newJSON *string
	if oldModel != nil {
		if b, err := json.Marshal(oldModel); err == nil {
			s := string(b)
			oldJSON = &s
		}
	}
	if newModel != nil {
		if b, err := json.Marshal(newModel); err == nil {
			s := string(b)
			newJSON = &s
		}
	}

	var userID *int64
	if v := tx.Statement.Context.Value("user_id"); v != nil {
		if uid, ok := v.(int64); ok {
			userID = &uid
		}
	}

	var ip *string
	if v := tx.Statement.Context.Value("ip"); v != nil {
		if s, ok := v.(string); ok {
			ip = &s
		}
	}

	var api *string
	if v := tx.Statement.Context.Value("api"); v != nil {
		if s, ok := v.(string); ok {
			api = &s
		}
	}

	history := History{
		UserID:   userID,
		Table:    &table,
		RowID:    &rowID,
		Action:   &action,
		OldValue: oldJSON,
		NewValue: newJSON,
		IP:       ip,
		API:      api,
	}

	_ = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&history)
}

// getPrimaryKey – modelning ID sini olish
func getPrimaryKey(model interface{}) interface{} {
	if model == nil {
		return nil
	}
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}
	field := v.FieldByName("ID")
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}
