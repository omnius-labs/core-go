package migration

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MySQLMigrator struct {
	db          *sqlx.DB
	path        string
	username    string
	description string
}

func NewMySQLMigrator(url string, path string, username string, description string) (*MySQLMigrator, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to database connect")
	}
	return &MySQLMigrator{db, path, username, description}, nil
}

func (m *MySQLMigrator) Migrate() error {
	err := m.init()
	if err != nil {
		return errors.Wrap(err, "failed to init")
	}

	histories, err := m.fetchMigrationHistories()
	if err != nil {
		return errors.Wrap(err, "failed to fetchMigrationHistories")
	}

	ignoreSet := make(map[string]bool)
	for _, history := range histories {
		ignoreSet[history.Filename] = true
	}

	files, err := m.loadMigrationFiles()
	if err != nil {
		return errors.Wrap(err, "failed to loadMigrationFiles")
	}

	var filteredFiles []*MySQLMigrationFile
	for _, file := range files {
		if !ignoreSet[file.Filename] {
			filteredFiles = append(filteredFiles, file)
		}
	}

	if len(filteredFiles) == 0 {
		return nil
	}

	err = m.semaphoreLock()
	if err != nil {
		return errors.Wrap(err, "failed to semaphoreLock")
	}

	result := m.executeMigrationQueries(filteredFiles)

	err = m.semaphoreUnlock()
	if err != nil {
		return errors.Wrap(err, "failed to semaphoreUnlock")
	}

	return result
}

func (m *MySQLMigrator) init() error {
	queries := `
CREATE TABLE IF NOT EXISTS _migrations (
    filename VARCHAR(255) NOT NULL,
    queries TEXT NOT NULL,
    executed_at DATETIME default CURRENT_TIMESTAMP,
    PRIMARY KEY (filename)
);
CREATE TABLE IF NOT EXISTS _semaphores (
    username VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    executed_at DATETIME default CURRENT_TIMESTAMP,
    PRIMARY KEY (username)
);
`
	return executeMultipleQueries(m.db, queries)
}

func (m *MySQLMigrator) fetchMigrationHistories() ([]*MySQLMigrationHistory, error) {
	var histories []*MySQLMigrationHistory
	err := m.db.Select(&histories, "SELECT filename, executed_at FROM _migrations")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return histories, nil
}

func (m *MySQLMigrator) loadMigrationFiles() ([]*MySQLMigrationFile, error) {
	var results []*MySQLMigrationFile

	entries, err := os.ReadDir(m.path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			content, err := os.ReadFile(filepath.Join(m.path, name))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			queries := string(content)
			result := &MySQLMigrationFile{Filename: name, Queries: queries}
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Filename < results[j].Filename
	})

	return results, nil
}

func (m *MySQLMigrator) executeMigrationQueries(files []*MySQLMigrationFile) error {
	for _, f := range files {
		tx, err := m.db.Beginx()
		if err != nil {
			return errors.WithStack(err)
		}

		err = executeMultipleQueriesTx(tx, f.Queries)
		if err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}

		err = m.insertMigrationHistory(tx, f.Filename, f.Queries)
		if err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}

		err = tx.Commit()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *MySQLMigrator) insertMigrationHistory(tx *sqlx.Tx, filename string, queries string) error {
	statement := "INSERT INTO _migrations (filename, queries) VALUES (?, ?)"
	_, err := tx.Exec(statement, filename, queries)
	return errors.WithStack(err)
}

func (m *MySQLMigrator) semaphoreLock() error {
	query := "INSERT INTO _semaphores (username, description) VALUES (?, ?)"
	_, err := m.db.Exec(query, m.username, m.description)
	return errors.WithStack(err)
}

func (m *MySQLMigrator) semaphoreUnlock() error {
	query := "DELETE FROM _semaphores WHERE username = ?"
	_, err := m.db.Exec(query, m.username)
	return errors.WithStack(err)
}

func executeMultipleQueries(db *sqlx.DB, queries string) error {
	splitQueries := strings.Split(queries, ";")

	for _, query := range splitQueries {
		query = strings.TrimSpace(query)
		if query != "" {
			_, err := db.Exec(query)
			if err != nil {
				return errors.Wrapf(err, "failed to execute query: %s", query)
			}
		}
	}

	return nil
}

func executeMultipleQueriesTx(tx *sqlx.Tx, queries string) error {
	splitQueries := strings.Split(queries, ";")

	for _, query := range splitQueries {
		query = strings.TrimSpace(query)
		if query != "" {
			_, err := tx.Exec(query)
			if err != nil {
				return errors.Wrapf(err, "failed to execute query: %s", query)
			}
		}
	}

	return nil
}

type MySQLMigrationFile struct {
	Filename string `db:"filename"`
	Queries  string `db:"queries"`
}

type MySQLMigrationHistory struct {
	Filename   string     `db:"filename"`
	ExecutedAt *time.Time `db:"executed_at"`
}
