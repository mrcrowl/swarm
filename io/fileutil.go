package io

import (
	"bufio"
	"io/ioutil"
	"os"
)

// ReadFirstLine reads the first line of a text file as a string
func ReadFirstLine(filepath string) (string, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		return sc.Text(), nil
	}

	return "", nil
}

// ReadContents reads the entire contents of a text file as a string
func ReadContents(filepath string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return string(bytes), nil
	}

	return "", err
}
