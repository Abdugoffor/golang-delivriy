package config

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func CreateHistoryTriggers(DB *gorm.DB, models []interface{}) {
	driver := DB.Dialector.Name()

	switch driver {

	// ================== 🟩 POSTGRES ==================
	case "postgres":
		createPostgresHistory(DB, models)

	// ================== 🟨 MYSQL ==================
	case "mysql":
		createMySQLHistory(DB, models)

	// ================== 🟦 SQLITE ==================
	case "sqlite":
		createSQLiteHistory(DB)

	default:
		log.Println("⚠️ Unsupported DB:", driver)
	}
}

func createPostgresHistory(DB *gorm.DB, models []interface{}) {

	DB.Exec(`
	CREATE TABLE IF NOT EXISTS histories (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NULL,
		table_name VARCHAR(100),
		row_id BIGINT NULL,
		action VARCHAR(20),
		ip VARCHAR(45),
		method VARCHAR(10),
		api TEXT,
		old_value JSONB,
		new_value JSONB,
		created_at TIMESTAMP DEFAULT now()
	);
	`)

	DB.Exec(`DROP FUNCTION IF EXISTS log_history CASCADE;`)

	DB.Exec(`
	CREATE FUNCTION log_history()
	RETURNS TRIGGER AS $$
	DECLARE
		v_row_id BIGINT;
	BEGIN
		IF (TG_OP = 'INSERT') THEN
			v_row_id := NEW.id;
			INSERT INTO histories(table_name,row_id,action,new_value)
			VALUES (TG_TABLE_NAME, v_row_id, 'INSERT', to_jsonb(NEW));
			RETURN NEW;

		ELSIF (TG_OP = 'UPDATE') THEN
			v_row_id := NEW.id;
			INSERT INTO histories(table_name,row_id,action,old_value,new_value)
			VALUES (TG_TABLE_NAME, v_row_id, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW));
			RETURN NEW;

		ELSIF (TG_OP = 'DELETE') THEN
			v_row_id := OLD.id;
			INSERT INTO histories(table_name,row_id,action,old_value)
			VALUES (TG_TABLE_NAME, v_row_id, 'DELETE', to_jsonb(OLD));
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;
	`)

	for _, model := range models {
		stmt := &gorm.Statement{DB: DB}
		_ = stmt.Parse(model)

		table := stmt.Schema.Table
		trigger := table + "_history"

		DB.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON %s;`, trigger, table))

		DB.Exec(fmt.Sprintf(`
		CREATE TRIGGER %s
		AFTER INSERT OR UPDATE OR DELETE ON %s
		FOR EACH ROW EXECUTE FUNCTION log_history();
		`, trigger, table))

		log.Println("✅ Postgres trigger:", table)
	}
}

func createMySQLHistory(DB *gorm.DB, models []interface{}) {

	DB.Exec(`
	CREATE TABLE IF NOT EXISTS histories (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		table_name VARCHAR(100),
		row_id BIGINT,
		action VARCHAR(20),
		old_value JSON,
		new_value JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)

	for _, model := range models {
		stmt := &gorm.Statement{DB: DB}
		_ = stmt.Parse(model)

		table := stmt.Schema.Table

		DB.Exec(fmt.Sprintf(`
		CREATE TRIGGER %s_insert AFTER INSERT ON %s
		FOR EACH ROW
		INSERT INTO histories(table_name,row_id,action)
		VALUES('%s', NEW.id, 'INSERT');
		`, table, table, table))

		log.Println("✅ MySQL trigger:", table)
	}
}

func createSQLiteHistory(DB *gorm.DB) {

	err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS histories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		table_name TEXT,
		row_id INTEGER,
		action TEXT,
		old_value TEXT,
		new_value TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`).Error

	if err != nil {
		log.Fatal("❌ Failed to create SQLite histories table:", err)
	}

	log.Println("✅ SQLite histories table ready")
}
