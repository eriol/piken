package sql // import "eriol.xyz/piken/sql"

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTablesQuery = `
	CREATE TABLE unicode_data (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT NOT NULL,
		canonical_class NUMERIC NOT NULL,
		bidi_class TEXT NOT NULL,
		decomposition_type TEXT,
		numeric_value_1 TEXT,
		numeric_value_2 TEXT,
		numeric_value_3 TEXT,
		bidi_mirrored TEXT,
		unicode_1_name TEXT,
		iso_comment TEXT,
		simple_uppercase_mapping TEXT,
		simple_lowercase_mapping TEXT,
		simple_titlecase_mapping TEXT
	);`
	insertUnicodeDataQuery = `INSERT INTO unicode_data (
		id,
		name,
		category,
		canonical_class, 
		bidi_class,
		decomposition_type,
		numeric_value_1,
		numeric_value_2,
		numeric_value_3,
		bidi_mirrored,
		unicode_1_name,
		iso_comment,
		simple_uppercase_mapping,
		simple_lowercase_mapping,
		simple_titlecase_mapping)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

type Store struct {
	db *sql.DB
}

// Open SQLite 3 database used by piken or create it if it doesn't exist yet.
func (s *Store) Open(database string) error {

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}

	s.db = db

	if _, err := os.Stat(database); err != nil {
		if os.IsNotExist(err) {
			if err := s.createDatabase(); err != nil {
				return err
			}
		}
	}

	return nil

}

// Close SQLite3 database used by piken.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) createDatabase() error {

	_, err := s.db.Exec(createTablesQuery)
	if err != nil {
		return err
	}
	return nil
}

// Load multiple records into piken database.
func (s *Store) LoadFromRecords(records [][]string) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertUnicodeDataQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {

		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}

		if _, err = stmt.Exec(args...); err != nil {
			return err
		}
	}
	tx.Commit()

	return nil
}
