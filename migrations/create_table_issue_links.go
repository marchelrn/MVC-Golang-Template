package migrations

import (
	"database/sql"
	"log"
)

type createIssueLinksTable struct{}

func (m *createIssueLinksTable) SkipProd() bool {
	return false
}

func getCreateIssueLinksTable() migration {
	return &createIssueLinksTable{}
}

func (m *createIssueLinksTable) Name() string {
	return "create-issue-links"
}

func (m *createIssueLinksTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.issue_links (
			id BIGSERIAL PRIMARY KEY,
			source_issue_id BIGINT NOT NULL REFERENCES issues(id),
			target_issue_id BIGINT NOT NULL REFERENCES issues(id),
			link_type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	log.Println("Creating Up migration: create-issue-links")
	return err
}

func (m *createIssueLinksTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS issue_links`)
	if err != nil {
		return err
	}
	return err
}

