package migrations

import (
	"database/sql"
	"log"
)

type createIssuesTable struct{}

func (m *createIssuesTable) SkipProd() bool {
	return false
}

func getCreateIssuesTable() migration {
	return &createIssuesTable{}
}

func (m *createIssuesTable) Name() string {
	return "create-issues"
}

func (m *createIssuesTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.issues (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id),
			sprint_id BIGINT REFERENCES sprints(id),
			epic_id BIGINT REFERENCES epics(id),
			type_id BIGINT NOT NULL REFERENCES issue_types(id),
			priority_id BIGINT NOT NULL REFERENCES priorities(id),
			status_id BIGINT NOT NULL REFERENCES statuses(id),
			reporter_id BIGINT NOT NULL REFERENCES users(id),
			assignee_id BIGINT REFERENCES users(id),
			parent_id BIGINT REFERENCES issues(id),
			issue_key VARCHAR(20) UNIQUE NOT NULL,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			story_points INT,
			original_estimate INT,
			time_spent INT DEFAULT 0,
			due_date DATE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_issues_project ON issues(project_id);
		CREATE INDEX IF NOT EXISTS idx_issues_sprint ON issues(sprint_id);
		CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee_id);
		CREATE INDEX IF NOT EXISTS idx_issues_reporter ON issues(reporter_id);
		CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status_id);
		CREATE INDEX IF NOT EXISTS idx_issues_key ON issues(issue_key);
		CREATE INDEX IF NOT EXISTS idx_issues_created ON issues(created_at);
	`)
	log.Println("Creating Up migration: create-issues")
	return err
}

func (m *createIssuesTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS issues`)
	if err != nil {
		return err
	}
	return err
}

