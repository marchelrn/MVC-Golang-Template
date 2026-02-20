package migrations

import (
	"database/sql"
	"log"
)

type createIssueLabelsTable struct{}

func (m *createIssueLabelsTable) SkipProd() bool {
	return false
}

func getCreateIssueLabelsTable() migration {
	return &createIssueLabelsTable{}
}

func (m *createIssueLabelsTable) Name() string {
	return "create-issue-labels"
}

func (m *createIssueLabelsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.issue_labels (
			id BIGSERIAL PRIMARY KEY,
			issue_id BIGINT NOT NULL REFERENCES issues(id),
			label_id BIGINT NOT NULL REFERENCES labels(id),
			UNIQUE (issue_id, label_id)
		);
	`)
	log.Println("Creating Up migration: create-issue-labels")
	return err
}

func (m *createIssueLabelsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS issue_labels`)
	if err != nil {
		return err
	}
	return err
}

