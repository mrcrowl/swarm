package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	s3Region            = "ap-southeast-2"
	s3Bucket            = "test.languageperfect.com"
	s3Path              = "/swarm"
	versionFilename     = "version.json"
	credentialsFilename = "c:\\wf\\lp\\Secrets\\credentials"
	credentialsProfile  = "default"
)

var swarmLocalPath = "../"

// JSON represents the version.json file from the remote server
type JSON struct {
	Version string `json:"version"`
}

func main() {
	cwd, _ := os.Getwd()
	if strings.HasSuffix(cwd, "swarm") {
		swarmLocalPath = "./"
	}

	ver := readver()
	localFilename, remoteFilename := filenames(ver)

	fmt.Printf("Target bucket: %s\n", s3Bucket)

	fmt.Println("Compressing...")
	gzippedBytes := gzipFile(localFilename)

	fmt.Printf("Writing to %s/%s...\n", s3Path, remoteFilename)
	writeFileToS3(gzippedBytes, remoteFilename, true)

	fmt.Printf("Updating %s/%s...\n", s3Path, versionFilename)
	versionBytes := encodeVersion(ver)
	writeFileToS3(versionBytes, versionFilename, false)
}

func encodeVersion(ver string) []byte {
	version := &JSON{Version: ver}
	bytes, _ := json.Marshal(version)
	return bytes
}

func readver() string {
	bytes, _ := ioutil.ReadFile(filepath.Join(swarmLocalPath, "main.go"))
	re := regexp.MustCompile("const localver = \"([^\"]+)\"")
	matches := re.FindSubmatch(bytes)
	if matches != nil {
		return string(matches[1])
	}

	panic("Could not read localver from main.go")
}

func filenames(ver string) (localFilename string, remoteFilename string) {
	suffix := ".exe"
	if runtime.GOOS != "windows" {
		suffix = ""
	}
	localFilename = fmt.Sprintf("swarm%s", suffix)
	remoteFilename = fmt.Sprintf("swarm-%s-%s%s", ver, runtime.GOOS, suffix)
	return
}

func gzipFile(localFilename string) []byte {
	// Get file size and read the file content into a buffer
	rawBytes, _ := ioutil.ReadFile(filepath.Join(swarmLocalPath, localFilename))

	var buffer bytes.Buffer
	gzwriter, _ := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	gzwriter.Write(rawBytes)
	gzippedBytes := buffer.Bytes()
	return gzippedBytes
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func writeFileToS3(gzippedBytes []byte, remoteFilename string, gzipped bool) {
	// Create a single AWS session (we can re use this if we're uploading many files)
	creds := credentials.NewEnvCredentials()
	if runtime.GOOS == "windows" {
		creds = credentials.NewSharedCredentials(credentialsFilename, credentialsProfile)
	}
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Region),
		Credentials: creds,
	})
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(s)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.

	reader := bytes.NewReader(gzippedBytes)
	key := path.Join(s3Path, remoteFilename)
	contentType := "application/octet-stream"
	cacheControl := ""
	if strings.HasSuffix(remoteFilename, ".json") {
		contentType = "application/json"
		cacheControl = "no-cache, no-store, must-revalidate"
	} else {
		// ensure executable file doesn't already exist
		headInput := &s3.HeadObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(key),
		}
		_, err := s3Client.HeadObject(headInput)
		if err == nil {
			panic(fmt.Sprintf("Have you updated localver in main.go?\n\nRemote file already exists: %s", remoteFilename))
		}
	}

	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(s3Bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(int64(len(gzippedBytes))),
		ContentType:   aws.String(contentType),
		CacheControl:  aws.String(cacheControl),
	}
	if gzipped {
		putInput.ContentEncoding = aws.String("gzip")
	}
	_, err = s3Client.PutObject(putInput)
	if err != nil {
		panic(err)
	}
}
