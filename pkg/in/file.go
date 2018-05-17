package in

import (
	"io"
	"io/ioutil"
	"os"
)

func ReadFile(filename string) (string, error) {

	var file io.Reader

	if filename == "-" {
		file = os.Stdin
	} else {
		reader, err := os.Open(filename)

		if err != nil {
			return "", err
		}

		file = reader
	}

	b, err := ioutil.ReadAll(file)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
