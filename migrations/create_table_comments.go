package migrations

import (
	"database/sql"
	"log"
)

type createCommentsTable struct{}

func (m *createCommentsTable) SkipProd() bool {
	return false
}

func getCreateCommentsTable() migration {
	return &createCommentsTable{}
}

func (m *createCommentsTable) Name() string {
	return "create-comments"
}

func (m *createCommentsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.comments (
			id BIGSERIAL PRIMARY KEY,
			issue_id BIGINT NOT NULL REFERENCES issues(id),
			user_id BIGINT NOT NULL REFERENCES users(id),
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_comments_issue ON comments(issue_id);
		CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);
	`)
	log.Println("Creating Up migration: create-comments")
	return err
}

func (m *createCommentsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS comments`)
	if err != nil {
		return err
	}
	return err
}

