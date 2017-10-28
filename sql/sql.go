// Copyright © 2015 Daniele Tricoli <eriol@mornie.org>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql // import "eriol.xyz/piken/sql"

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// Full-text search need only TEXT column.
	createTablesQuery = `
	CREATE TABLE last_update (
		filename TEXT NOT NULL PRIMARY KEY,
		date TEXT
	);
	CREATE VIRTUAL TABLE unicode_data USING fts4(
		codepoint TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT NOT NULL,
		canonical_class TEXT NOT NULL,
		bidi_class TEXT NOT NULL,
		decomposition_type TEXT,
		numeric_type TEXT,
		numeric_digit TEXT,
		numeric_value TEXT,
		bidi_mirrored TEXT,
		unicode_1_name TEXT,
		iso_comment TEXT,
		simple_uppercase_mapping TEXT,
		simple_lowercase_mapping TEXT,
		simple_titlecase_mapping TEXT
	);`
	insertUnicodeDataQuery = `INSERT INTO unicode_data (
		codepoint,
		name,
		category,
		canonical_class, 
		bidi_class,
		decomposition_type,
		numeric_type,
		numeric_digit,
		numeric_value,
		bidi_mirrored,
		unicode_1_name,
		iso_comment,
		simple_uppercase_mapping,
		simple_lowercase_mapping,
		simple_titlecase_mapping)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	deleteUdicodeData     = `DELETE FROM unicode_data;`
	createLastUpdateQuery = `INSERT INTO last_update (filename, date) VALUES (?, ?)`
	updateLastUpdateQuery = `UPDATE last_update SET date = ? where filename = ?`
	getLastUpdateQuery    = `SELECT date FROM last_update WHERE filename = ?`
	getUnicodeQuery       = `SELECT * FROM unicode_data WHERE name MATCH ?`
	countUnicodeQuery     = `SELECT count(*) FROM unicode_data`
)

type Store struct {
	db *sql.DB
}

type UnicodeData struct {
	CodePoint              string
	Name                   string
	Category               string
	CanonicalClass         string
	BidiClass              string
	DecompositionType      string
	NumericType            string
	NumericDigit           string
	NumericValue           string
	BidiMirrored           string
	Unicode1Name           string
	ISOComment             string
	SimpleUppercaseMapping string
	SimpleLowercaseMapping string
	SimpleTitlecaseMapping string
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

// DeleteUnicodeData removes all data from unicode_data table.
func (s *Store) DeleteUnicodeData() error {

	_, err := s.db.Exec(deleteUdicodeData)
	if err != nil {
		return err
	}
	return nil
}

// Create latest update entry for given file.
func (s *Store) CreateLastUpdate(filename string, t time.Time) error {

	ts := t.Format(time.RFC3339)

	_, err := s.db.Exec(createLastUpdateQuery, filename, ts)

	return err
}

// UpdateLastUpdate updates date entry for given file.
func (s *Store) UpdateLastUpdate(filename string, t time.Time) error {

	ts := t.Format(time.RFC3339)

	_, err := s.db.Exec(updateLastUpdateQuery, ts, filename)

	return err
}

// Get latest update for given file.
func (s *Store) GetLastUpdate(filename string) (time.Time, error) {
	var t string

	s.db.QueryRow(getLastUpdateQuery, filename).Scan(&t)

	// If the query is empty return the beginning of time.
	if t == "" {
		return time.Unix(0, 0), nil
	}

	tp, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return time.Unix(0, 0), err
	}

	return tp, nil
}

// Search unicode data using name.
func (s *Store) SearchUnicode(name string) (records []UnicodeData, err error) {

	rows, err := s.db.Query(getUnicodeQuery, name)
	defer rows.Close()
	if err != nil {
		return []UnicodeData{}, err
	}

	for rows.Next() {
		var row UnicodeData

		if err := rows.Scan(
			&row.CodePoint,
			&row.Name,
			&row.Category,
			&row.CanonicalClass,
			&row.BidiClass,
			&row.DecompositionType,
			&row.NumericType,
			&row.NumericDigit,
			&row.NumericValue,
			&row.BidiMirrored,
			&row.Unicode1Name,
			&row.ISOComment,
			&row.SimpleUppercaseMapping,
			&row.SimpleLowercaseMapping,
			&row.SimpleTitlecaseMapping); err != nil {
			return []UnicodeData{}, err
		}

		records = append(records, row)

	}

	return records, nil
}

// Count total rows inside unicode data table.
func (s *Store) CountUnicodeData() int {
	var total int

	s.db.QueryRow(countUnicodeQuery).Scan(&total)

	return total
}
