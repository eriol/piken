package sql

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Insert some data for testing.
func fixture(databaseName string) {
	db, _ := sql.Open("sqlite3", databaseName)
	defer db.Close()

	db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F60E', 'SMILING FACE WITH SUNGLASSES', 'So')`)
	db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F602', 'FACE WITH TEARS OF JOY', 'So')`)
	db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F4D7', 'GREEN BOOK', 'So')`)
}

func TestSearch(t *testing.T) {
	var s Store

	dirName, err := ioutil.TempDir("", "piken")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)

	databaseName := path.Join(dirName, "piken.sqlite3")

	if err := s.Open(databaseName); err != nil {
		assert.Error(t, err)
	}
	defer s.Close()

	fixture(databaseName)

	data, err := s.SearchUnicode("cat")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, data, []UnicodeData(nil))

	data, err = s.SearchUnicode("book")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].CodePoint, "1F4D7")
	assert.Equal(t, data[0].Name, "GREEN BOOK")
	assert.Equal(t, data[0].Category, "So")

	// Search is case insensitive
	data, err = s.SearchUnicode("fAcE")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, len(data), 2)
	assert.Equal(t, data[0].CodePoint, "1F60E")
	assert.Equal(t, data[0].Name, "SMILING FACE WITH SUNGLASSES")
	assert.Equal(t, data[0].Category, "So")
	assert.Equal(t, data[1].CodePoint, "1F602")
	assert.Equal(t, data[1].Name, "FACE WITH TEARS OF JOY")
	assert.Equal(t, data[1].Category, "So")

	data, err = s.SearchUnicode("face NOT sunglasses")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].CodePoint, "1F602")
	assert.Equal(t, data[0].Name, "FACE WITH TEARS OF JOY")
	assert.Equal(t, data[0].Category, "So")
}
