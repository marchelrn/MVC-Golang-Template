package migrations

import (
	"database/sql"
	"log"
)

type createIssueHistoryTable struct{}

func (m *createIssueHistoryTable) SkipProd() bool {
	return false
}

func getCreateIssueHistoryTable() migration {
	return &createIssueHistoryTable{}
}

func (m *createIssueHistoryTable) Name() string {
	return "create-issue-history"
}

func (m *createIssueHistoryTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.issue_history (
			id BIGSERIAL PRIMARY KEY,
			issue_id BIGINT NOT NULL REFERENCES issues(id),
			user_id BIGINT NOT NULL REFERENCES users(id),
			field_name VARCHAR(100) NOT NULL,
			old_value TEXT,
			new_value TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_issue_history_issue ON issue_history(issue_id);
		CREATE INDEX IF NOT EXISTS idx_issue_history_created ON issue_history(created_at);
	`)
	log.Println("Creating Up migration: create-issue-history")
	return err
}

func (m *createIssueHistoryTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS issue_history`)
	if err != nil {
		return err
	}
	return err
}

