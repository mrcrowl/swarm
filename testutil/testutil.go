package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// RemoveTempDir creates a temporary directory
func RemoveTempDir(tempDir string) {
	if strings.HasPrefix(tempDir, os.TempDir()) {
		os.RemoveAll(tempDir)
	} else {
		panic(fmt.Sprintf("RemoveTempDir was asked to delete non-temp path!!!"))
	}
}

// MakeSubdirectoryTree creates necessary folders for a path to exist
func MakeSubdirectoryTree(parentPath string, subdirectoryPath string) string {
	targetPath := filepath.Join(parentPath, subdirectoryPath)
	err := os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("MakeSubdirectoryTree for '%s' failed: %s", targetPath, err))
	}
	return targetPath
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
