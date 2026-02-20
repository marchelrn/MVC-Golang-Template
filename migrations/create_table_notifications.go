package migrations

import (
	"database/sql"
	"log"
)

type createNotificationsTable struct{}

func (m *createNotificationsTable) SkipProd() bool {
	return false
}

func getCreateNotificationsTable() migration {
	return &createNotificationsTable{}
}

func (m *createNotificationsTable) Name() string {
	return "create-notifications"
}

func (m *createNotificationsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.notifications (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id),
			issue_id BIGINT REFERENCES issues(id),
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
		CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(user_id, is_read);
	`)
	log.Println("Creating Up migration: create-notifications")
	return err
}

func (m *createNotificationsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS notifications`)
	if err != nil {
		return err
	}
	return err
}

