package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateTempDirWithPrefix creates a temporary directory
func CreateTempDirWithPrefix(prefix string) string {
	temppath, _ := ioutil.TempDir("", prefix)
	return temppath
}

// CreateTempDir creates a temporary directory
func CreateTempDir() string {
	temppath, _ := ioutil.TempDir("", "swarm-temp")
	return temppath
}

// WriteTextFile writes a string as file
func WriteTextFile(folderPath string, filename string, contents string) string {
	absoluteFilepath := filepath.Join(folderPath, filename)
	bytes := []byte(contents)
	ioutil.WriteFile(absoluteFilepath, bytes, os.ModePerm)
	return absoluteFilepath
}

// ReadTextFile reads a file as a string
func ReadTextFile(folderPath string, filename string) string {
	absoluteFilepath := filepath.Join(folderPath, filename)
	bytes, _ := ioutil.ReadFile(absoluteFilepath)
	text := string(bytes)
	return text
}
