package sql

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var store Store

func (s *Store) init() (dirName string, err error) {

	dirName, err = ioutil.TempDir("", "piken")
	if err != nil {
		return "", err
	}

	databaseName := path.Join(dirName, "piken.sqlite3")

	if err := s.Open(databaseName); err != nil {
		return "", err
	}

	return dirName, nil
}

// Insert some data for testing.
func (s *Store) fixture() {
	s.db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F60E', 'SMILING FACE WITH SUNGLASSES', 'So')`)
	s.db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F602', 'FACE WITH TEARS OF JOY', 'So')`)
	s.db.Exec(`INSERT INTO unicode_data (codepoint, name, category)
		VALUES ('1F4D7', 'GREEN BOOK', 'So')`)
}

func TestSearch(t *testing.T) {

	dirName, err := store.init()
	if err != nil {
		assert.Error(t, err)
	}
	defer os.RemoveAll(dirName)
	defer store.Close()

	store.fixture()

	data, err := store.SearchUnicode("cat")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, data, []UnicodeData(nil))

	data, err = store.SearchUnicode("book")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].CodePoint, "1F4D7")
	assert.Equal(t, data[0].Name, "GREEN BOOK")
	assert.Equal(t, data[0].Category, "So")

	// Search is case insensitive
	data, err = store.SearchUnicode("fAcE")
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

	data, err = store.SearchUnicode("face NOT sunglasses")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].CodePoint, "1F602")
	assert.Equal(t, data[0].Name, "FACE WITH TEARS OF JOY")
	assert.Equal(t, data[0].Category, "So")
}

func TestLastUpdate(t *testing.T) {
	dirName, err := store.init()
	if err != nil {
		assert.Error(t, err)
	}
	defer os.RemoveAll(dirName)
	defer store.Close()

	updateTime, err := store.GetLastUpdate("test.txt")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, updateTime, time.Unix(0, 0))

	testTime, _ := time.Parse("2006-Jan-02", "1993-Jan-25")
	if err := store.CreateLastUpdate("test.txt", testTime); err != nil {
		assert.Error(t, err)
	}
	updateTime, err = store.GetLastUpdate("test.txt")
	assert.Equal(t, updateTime, testTime)
}

func TestLoadFromRecords(t *testing.T) {
	dirName, err := store.init()
	if err != nil {
		assert.Error(t, err)
	}
	defer os.RemoveAll(dirName)
	defer store.Close()

	records := [][]string{
		[]string{"1F60E", "SMILING FACE WITH SUNGLASSES", "So", "", "", "", "", "", "", "", "", "", "", "", ""},
		[]string{"1F602", "FACE WITH TEARS OF JOY", "So", "", "", "", "", "", "", "", "", "", "", "", ""},
		[]string{"1F4D7", "GREEN BOOK", "So", "", "", "", "", "", "", "", "", "", "", "", ""},
	}

	store.LoadFromRecords(records)

	assert.Equal(t, store.CountUnicodeData(), 3)
}
