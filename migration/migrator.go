package migration

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Migrator struct {
	db          *sql.DB
	path        string
	username    string
	description string
}

func NewMigrator(url string, path string, username string, description string) (*Migrator, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Migrator{db, path, username, description}, nil
}

func (m *Migrator) Migrate() error {
	m.init()
	return nil
}

func (m *Migrator) init() error {
	queries := `
create table if not exists _migrations (
    filename varchar(255) NOT NULL,
    queries text NOT NULL,
    executed_at timestamp without time zone default CURRENT_TIMESTAMP,
    primary key (filename)
);
create table if not exists _semaphores (
    username varchar(255) NOT NULL,
    description text NOT NULL,
    executed_at timestamp without time zone default CURRENT_TIMESTAMP,
    primary key (username)
);
`
	m.db.Exec(queries)
	return nil
}
