package migrations

import (
	"database/sql"
	"log"
)

type createAttachmentsTable struct{}

func (m *createAttachmentsTable) SkipProd() bool {
	return false
}

func getCreateAttachmentsTable() migration {
	return &createAttachmentsTable{}
}

func (m *createAttachmentsTable) Name() string {
	return "create-attachments"
}

func (m *createAttachmentsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.attachments (
			id BIGSERIAL PRIMARY KEY,
			issue_id BIGINT NOT NULL REFERENCES issues(id),
			uploaded_by BIGINT NOT NULL REFERENCES users(id),
			filename VARCHAR(255) NOT NULL,
			filepath VARCHAR(500) NOT NULL,
			mimetype VARCHAR(100) NOT NULL,
			filesize BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	log.Println("Creating Up migration: create-attachments")
	return err
}

func (m *createAttachmentsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS attachments`)
	if err != nil {
		return err
	}
	return err
}

