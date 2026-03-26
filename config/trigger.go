package config

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func CreateHistoryTriggers(DB *gorm.DB, models []interface{}) {
	driver := strings.ToLower(DB.Dialector.Name())

	switch driver {
	case "postgres", "postgresql":
		createPostgresHistory(DB, models)

	case "mysql":
		createMySQLHistory(DB, models)

	case "sqlite", "sqlite3":
		createSQLiteHistory(DB, models)

	default:
		log.Printf("⚠️ Unsupported DB for history triggers: %s", driver)
	}
}

// ================== 🟩 POSTGRES ==================
func createPostgresHistory(DB *gorm.DB, models []interface{}) {

	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS histories (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NULL,
			table_name VARCHAR(100),
			row_id BIGINT NULL,
			action VARCHAR(20),
			ip VARCHAR(45) NULL,
			method VARCHAR(10) NULL,
			api TEXT NULL,
			old_value JSONB,
			new_value JSONB,
			created_at TIMESTAMP DEFAULT now()
		);
	`).Error; err != nil {
		log.Fatal("❌ Failed to create histories table: ", err)
	}

	if err := DB.Exec(`
		DO $$
		BEGIN
			CREATE OR REPLACE FUNCTION log_history()
			RETURNS TRIGGER AS $func$
			DECLARE
				v_old JSONB := '{}'::JSONB;
				v_new JSONB := '{}'::JSONB;
				v_uid_text TEXT;
				v_uid BIGINT;
				v_ip TEXT;
				v_path TEXT;
				v_method TEXT;
				v_row_id BIGINT;
				key TEXT;
				old_row JSONB;
				new_row JSONB;
			BEGIN
				v_uid_text := current_setting('app.current_user_id', true);
				IF v_uid_text IS NULL OR v_uid_text = '' THEN
					v_uid := NULL;
				ELSE
					v_uid := v_uid_text::BIGINT;
				END IF;

				v_ip := current_setting('app.current_request_ip', true);
				v_path := current_setting('app.current_request_path', true);
				v_method := current_setting('app.current_request_method', true);

				IF (TG_OP = 'INSERT') THEN
					v_new := to_jsonb(NEW);
					BEGIN v_row_id := NEW.id; EXCEPTION WHEN others THEN v_row_id := NULL; END;
					INSERT INTO histories (user_id, table_name, row_id, action, new_value, ip, api, method, created_at)
					VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'INSERT', v_new, v_ip, v_path, v_method, now());
					RETURN NEW;

				ELSIF (TG_OP = 'UPDATE') THEN
					old_row := to_jsonb(OLD);
					new_row := to_jsonb(NEW);

						-- faqat o‘zgargan maydonlarni yozamiz
					FOR key IN SELECT jsonb_object_keys(new_row) LOOP
						IF old_row->>key IS DISTINCT FROM new_row->>key THEN
							v_old := v_old || jsonb_build_object(key, old_row->key);
							v_new := v_new || jsonb_build_object(key, new_row->key);
						END IF;
					END LOOP;

					BEGIN v_row_id := NEW.id; EXCEPTION WHEN others THEN v_row_id := NULL; END;
					INSERT INTO histories (user_id, table_name, row_id, action, old_value, new_value, ip, api, method, created_at)
					VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'UPDATE', v_old, v_new, v_ip, v_path, v_method, now());
					RETURN NEW;

				ELSIF (TG_OP = 'DELETE') THEN
					v_old := to_jsonb(OLD);
					BEGIN v_row_id := OLD.id; EXCEPTION WHEN others THEN v_row_id := NULL; END;
					INSERT INTO histories (user_id, table_name, row_id, action, old_value, ip, api, method, created_at)
					VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'DELETE', v_old, v_ip, v_path, v_method, now());
					RETURN OLD;
				END IF;

				RETURN NULL;
			END;
			$func$ LANGUAGE plpgsql;
		END $$;
	`).Error; err != nil {
		log.Fatal("❌ Failed to create or replace log_history function: ", err)
	}

	for _, model := range models {
		stmt := &gorm.Statement{DB: DB}
		if err := stmt.Parse(model); err != nil {
			log.Fatalf("❌ Failed to parse model: %v", err)
		}

		table := stmt.Schema.Table
		if strings.TrimSpace(table) == "" {
			log.Printf("⚠️ Skip: model has no table name")
			continue
		}

		triggerName := fmt.Sprintf("%s_history", table)

		sql := fmt.Sprintf(`
			DO $$
			BEGIN
					-- eski triggerni o‘chiramiz
				IF EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = '%s') THEN
					EXECUTE format('DROP TRIGGER IF EXISTS %s ON %s;', '%s', '%s');
				END IF;

				-- yangisini yaratamiz
				EXECUTE format('CREATE TRIGGER %s AFTER INSERT OR UPDATE OR DELETE ON %s FOR EACH ROW EXECUTE FUNCTION log_history();', '%s', '%s');
			END $$;
		`, triggerName, triggerName, table, triggerName, table, triggerName, table)

		if err := DB.Exec(sql).Error; err != nil {
			log.Printf("❌ Failed to (re)create trigger for %s: %v", table, err)
		} else {
			log.Printf("✅ Trigger updated for table: %s", table)
		}
	}
}

// ================== 🟨 MYSQL ==================
func createMySQLHistory(DB *gorm.DB, models []interface{}) {
	DB.Exec(`
		CREATE TABLE IF NOT EXISTS histories (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NULL,
			table_name VARCHAR(100),
			row_id BIGINT NULL,
			action VARCHAR(20),
			ip VARCHAR(45) NULL,
			method VARCHAR(10) NULL,
			api TEXT NULL,
			old_value JSON,
			new_value JSON,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	for _, model := range models {
		stmt := &gorm.Statement{DB: DB}
		if err := stmt.Parse(model); err != nil {
			log.Printf("❌ Failed to parse model: %v", err)
			continue
		}

		table := stmt.Schema.Table
		if table == "" {
			continue
		}

		// Ustunlar ro'yxatini olish
		var columns []string
		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				columns = append(columns, field.DBName)
			}
		}

		if len(columns) == 0 {
			log.Printf("⚠️ No columns found for table: %s", table)
			continue
		}

		// JSON_OBJECT('col', NEW.col, ...) qurish
		newParts := make([]string, 0, len(columns))
		oldParts := make([]string, 0, len(columns))
		for _, col := range columns {
			newParts = append(newParts, fmt.Sprintf("'%s', NEW.%s", col, col))
			oldParts = append(oldParts, fmt.Sprintf("'%s', OLD.%s", col, col))
		}
		newJSON := fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(newParts, ", "))
		oldJSON := fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(oldParts, ", "))

		// Eski triggerlarni o'chirish
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_insert;", table))
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_update;", table))
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_delete;", table))

		// INSERT trigger
		insertSQL := fmt.Sprintf(`
			CREATE TRIGGER %s_insert AFTER INSERT ON %s
			FOR EACH ROW
			BEGIN
				INSERT INTO histories (user_id, table_name, row_id, action, new_value, ip, method, api, created_at)
				VALUES (@current_user_id, '%s', NEW.id, 'INSERT', %s,
						@current_request_ip, @current_request_method, @current_request_path, NOW());
			END;
		`, table, table, table, newJSON)

		// UPDATE trigger
		updateSQL := fmt.Sprintf(`
			CREATE TRIGGER %s_update AFTER UPDATE ON %s
			FOR EACH ROW
			BEGIN
				INSERT INTO histories (user_id, table_name, row_id, action, old_value, new_value, ip, method, api, created_at)
				VALUES (@current_user_id, '%s', NEW.id, 'UPDATE',
						%s, %s,
						@current_request_ip, @current_request_method, @current_request_path, NOW());
			END;
		`, table, table, table, oldJSON, newJSON)

		// DELETE trigger
		deleteSQL := fmt.Sprintf(`
			CREATE TRIGGER %s_delete AFTER DELETE ON %s
			FOR EACH ROW
			BEGIN
				INSERT INTO histories (user_id, table_name, row_id, action, old_value, ip, method, api, created_at)
				VALUES (@current_user_id, '%s', OLD.id, 'DELETE',
						%s,
						@current_request_ip, @current_request_method, @current_request_path, NOW());
			END;
		`, table, table, table, oldJSON)

		if err := DB.Exec(insertSQL).Error; err != nil {
			log.Printf("❌ MySQL INSERT trigger error for %s: %v", table, err)
		}
		if err := DB.Exec(updateSQL).Error; err != nil {
			log.Printf("❌ MySQL UPDATE trigger error for %s: %v", table, err)
		}
		if err := DB.Exec(deleteSQL).Error; err != nil {
			log.Printf("❌ MySQL DELETE trigger error for %s: %v", table, err)
		}

		log.Printf("✅ MySQL triggers created for: %s", table)
	}
}

// ================== 🟦 SQLITE ==================
func createSQLiteHistory(DB *gorm.DB, models []interface{}) {
	err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS histories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NULL,
			table_name TEXT,
			row_id INTEGER NULL,
			action TEXT,
			ip TEXT NULL,
			method TEXT NULL,
			api TEXT NULL,
			old_value TEXT,
			new_value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`).Error

	if err != nil {
		log.Fatal("❌ Failed to create SQLite histories table:", err)
	}

	for _, model := range models {
		stmt := &gorm.Statement{DB: DB}
		if err := stmt.Parse(model); err != nil {
			log.Printf("❌ Failed to parse model: %v", err)
			continue
		}

		table := stmt.Schema.Table
		if table == "" {
			continue
		}

		// Ustunlar ro'yxatini olish
		var columns []string
		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				columns = append(columns, field.DBName)
			}
		}

		if len(columns) == 0 {
			log.Printf("⚠️ No columns found for table: %s", table)
			continue
		}

		// json_object('col', NEW.col, ...) qurish
		newParts := make([]string, 0, len(columns))
		oldParts := make([]string, 0, len(columns))
		for _, col := range columns {
			newParts = append(newParts, fmt.Sprintf("'%s', NEW.%s", col, col))
			oldParts = append(oldParts, fmt.Sprintf("'%s', OLD.%s", col, col))
		}
		newJSON := fmt.Sprintf("json_object(%s)", strings.Join(newParts, ", "))
		oldJSON := fmt.Sprintf("json_object(%s)", strings.Join(oldParts, ", "))

		prefix := table + "_history"

		// Eski triggerlarni tozalash
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_insert;", prefix))
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_update;", prefix))
		DB.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s_delete;", prefix))

		// INSERT trigger
		insertSQL := fmt.Sprintf(`
			CREATE TRIGGER IF NOT EXISTS %s_insert AFTER INSERT ON %s
			BEGIN
				INSERT INTO histories (table_name, row_id, action, new_value, created_at)
				VALUES ('%s', NEW.id, 'INSERT', %s, CURRENT_TIMESTAMP);
			END;
		`, prefix, table, table, newJSON)

		// UPDATE trigger
		updateSQL := fmt.Sprintf(`
			CREATE TRIGGER IF NOT EXISTS %s_update AFTER UPDATE ON %s
			BEGIN
				INSERT INTO histories (table_name, row_id, action, old_value, new_value, created_at)
				VALUES ('%s', NEW.id, 'UPDATE', %s, %s, CURRENT_TIMESTAMP);
			END;
		`, prefix, table, table, oldJSON, newJSON)

		// DELETE trigger
		deleteSQL := fmt.Sprintf(`
			CREATE TRIGGER IF NOT EXISTS %s_delete AFTER DELETE ON %s
			BEGIN
				INSERT INTO histories (table_name, row_id, action, old_value, created_at)
				VALUES ('%s', OLD.id, 'DELETE', %s, CURRENT_TIMESTAMP);
			END;
		`, prefix, table, table, oldJSON)

		if err := DB.Exec(insertSQL).Error; err != nil {
			log.Printf("❌ SQLite INSERT trigger error for %s: %v", table, err)
		}
		if err := DB.Exec(updateSQL).Error; err != nil {
			log.Printf("❌ SQLite UPDATE trigger error for %s: %v", table, err)
		}
		if err := DB.Exec(deleteSQL).Error; err != nil {
			log.Printf("❌ SQLite DELETE trigger error for %s: %v", table, err)
		}

		log.Printf("✅ SQLite triggers created for table: %s", table)
	}

	log.Println("✅ SQLite history triggers setup completed")
}
