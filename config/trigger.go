package config

import (
	"fmt"
	"log"
)

func CreateHistoryTriggers(tables []string) {
	// 1. histories jadvali yaratiladi
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS histories (
			id SERIAL PRIMARY KEY,
			user_id BIGINT,
			table_name TEXT NOT NULL,
			row_id TEXT,
			action TEXT NOT NULL,
			old_value JSONB,
			new_value JSONB,
			ip TEXT,
			api TEXT,
			method TEXT,
			created_at TIMESTAMP DEFAULT now()
		);
	`).Error; err != nil {
		log.Fatal("❌ Failed to create histories table: ", err)
	}

	// 2. log_history function yo‘q bo‘lsa yaratamiz
	if err := DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_proc WHERE proname = 'log_history'
			) THEN
				CREATE FUNCTION log_history()
				RETURNS TRIGGER AS $func$
				DECLARE
					v_old JSONB;
					v_new JSONB;
					v_uid_text TEXT;
					v_uid BIGINT;
					v_ip TEXT;
					v_path TEXT;
					v_method TEXT;
					v_row_id TEXT;
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
						BEGIN v_row_id := NEW.id::text; EXCEPTION WHEN others THEN v_row_id := NULL; END;
						INSERT INTO histories (user_id, table_name, row_id, action, new_value, ip, api, method, created_at)
						VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'INSERT', v_new, v_ip, v_path, v_method, now());
						RETURN NEW;
					ELSIF (TG_OP = 'UPDATE') THEN
						v_old := to_jsonb(OLD); v_new := to_jsonb(NEW);
						BEGIN v_row_id := NEW.id::text; EXCEPTION WHEN others THEN v_row_id := NULL; END;
						INSERT INTO histories (user_id, table_name, row_id, action, old_value, new_value, ip, api, method, created_at)
						VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'UPDATE', v_old, v_new, v_ip, v_path, v_method, now());
						RETURN NEW;
					ELSIF (TG_OP = 'DELETE') THEN
						v_old := to_jsonb(OLD);
						BEGIN v_row_id := OLD.id::text; EXCEPTION WHEN others THEN v_row_id := NULL; END;
						INSERT INTO histories (user_id, table_name, row_id, action, old_value, ip, api, method, created_at)
						VALUES (v_uid, TG_TABLE_NAME, v_row_id, 'DELETE', v_old, v_ip, v_path, v_method, now());
						RETURN OLD;
					END IF;

					RETURN NULL;
				END;
				$func$ LANGUAGE plpgsql;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Fatal("❌ Failed to create log_history function: ", err)
	}

	// 3. Har bir jadval uchun trigger borligini tekshirish va qo‘shish
	for _, table := range tables {
		triggerName := fmt.Sprintf("%s_history", table)

		sql := fmt.Sprintf(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_trigger WHERE tgname = '%s'
				) THEN
					EXECUTE $q$
						CREATE TRIGGER %s
						AFTER INSERT OR UPDATE OR DELETE ON %s
						FOR EACH ROW EXECUTE FUNCTION log_history()
					$q$;
				END IF;
			END $$;
		`, triggerName, triggerName, table)

		if err := DB.Exec(sql).Error; err != nil {
			log.Fatalf("❌ Failed to create trigger for %s: %v", table, err)
		} else {
			log.Printf("✅ Trigger checked/created for table: %s", table)
		}
	}
}
