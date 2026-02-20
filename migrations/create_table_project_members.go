package migrations

import (
	"database/sql"
	"log"
)

type createProjectMembersTable struct{}

func (m *createProjectMembersTable) SkipProd() bool {
	return false
}

func getCreateProjectMembersTable() migration {
	return &createProjectMembersTable{}
}

func (m *createProjectMembersTable) Name() string {
	return "create-project-members"
}

func (m *createProjectMembersTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.project_members (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id),
			user_id BIGINT NOT NULL REFERENCES users(id),
			role VARCHAR(50) NOT NULL,
			joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE (project_id, user_id)
		);
	`)
	log.Println("Creating Up migration: create-project-members")
	return err
}

func (m *createProjectMembersTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS project_members`)
	if err != nil {
		return err
	}
	return err
}

