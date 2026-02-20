package migrations

import (
	"database/sql"
	"log"
)

type createIssueTypesTable struct{}

func (m *createIssueTypesTable) SkipProd() bool {
	return false
}

func getCreateIssueTypesTable() migration {
	return &createIssueTypesTable{}
}

func (m *createIssueTypesTable) Name() string {
	return "create-issue-types"
}

func (m *createIssueTypesTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.issue_types (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			description VARCHAR(255),
			icon VARCHAR(100),
			is_subtask BOOLEAN DEFAULT FALSE
		);
	`)
	log.Println("Creating Up migration: create-issue-types")
	return err
}

func (m *createIssueTypesTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS issue_types`)
	if err != nil {
		return err
	}
	return err
}

