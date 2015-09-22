package main // import "eriol.xyz/piken"

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
)

// Get Last-Modified for resources with the given URL.
func checkLastModified(url string) (time.Time, error) {
	r, err := http.Head(url)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer r.Body.Close()

	t, err := time.Parse(time.RFC1123, r.Header["Last-Modified"][0])
	if err != nil {
		return time.Unix(0, 0), err
	}

	return t, nil

}

func download(url, output string) error {

	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	// io.copyBuffer, the actual implementation of io.Copy, reads maximum 32 KB
	// from input, writes to output and then repeats. No need to worry about
	// the size of file to download.
	_, err = io.Copy(out, r.Body)
	if err != nil {
		return err
	}

	return nil

}

// Get user home directory or exit with a fatal error.
func getHome() string {

	homeDir, err := homedir.Dir()
	if err != nil {
		logrus.Fatal(err)
	}

	return homeDir
}

// Read a CSV file and return a slice of slice.
func readCsvFile(filepath string) (records [][]string, err error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err = reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil

}

// Convert an unicode codepoint into a string.
func codePointToGlyph(codepoint string) (string, error) {

	s, err := strconv.ParseInt(codepoint, 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%c", s), nil

}
