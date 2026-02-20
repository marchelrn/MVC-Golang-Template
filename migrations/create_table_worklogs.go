package migrations

import (
	"database/sql"
	"log"
)

type createWorklogsTable struct{}

func (m *createWorklogsTable) SkipProd() bool {
	return false
}

func getCreateWorklogsTable() migration {
	return &createWorklogsTable{}
}

func (m *createWorklogsTable) Name() string {
	return "create-worklogs"
}

func (m *createWorklogsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.worklogs (
			id BIGSERIAL PRIMARY KEY,
			issue_id BIGINT NOT NULL REFERENCES issues(id),
			user_id BIGINT NOT NULL REFERENCES users(id),
			time_spent INT NOT NULL,
			description TEXT,
			work_date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_worklogs_issue ON worklogs(issue_id);
		CREATE INDEX IF NOT EXISTS idx_worklogs_user ON worklogs(user_id);
	`)
	log.Println("Creating Up migration: create-worklogs")
	return err
}

func (m *createWorklogsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS worklogs`)
	if err != nil {
		return err
	}
	return err
}

