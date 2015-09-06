package main // import "eriol.xyz/piken"

import (
	"io"
	"net/http"
	"os"
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
