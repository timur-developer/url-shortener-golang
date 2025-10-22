package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(32) NOT NULL UNIQUE,
			password VARCHAR(64) NOT NULL,
	    	user_role VARCHAR(16) NOT NULL
		)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to create table: %w", op, err)
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, password, user_role)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO NOTHING
`, 0, "anonymous", "anonymous", "anonymous")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to insert anonymous user: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_alias ON users (username)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to create index: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL,
			url TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	    	UNIQUE (alias, user_id)
		)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to create table: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_alias ON url (alias)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to create index: %w", op, err)
	}

	return &Storage{db: db}, nil
}
