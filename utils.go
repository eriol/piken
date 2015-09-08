package main // import "eriol.xyz/piken"

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
)

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
