package migrations

import (
	"database/sql"
	"log"
)

type createEmailsVerifTable struct{}

func (m *createEmailsVerifTable) SkipProd() bool {
	return false
}

func getCreateEmailVerif() migration {
	return &createEmailsVerifTable{}
}

func (m *createEmailsVerifTable) Name() string {
	return "create-emails-verification"
}

func (m *createEmailsVerifTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.emails_verification (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email VARCHAR(255) NOT NULL,
			token VARCHAR(255) NOT NULL UNIQUE,
			expired_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT fk_emails_verif_user FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	log.Println("Creating Up migration: create-emails-verification")
	return err
}

func (m *createEmailsVerifTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS emails_verification`)
	return err
}
