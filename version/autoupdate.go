package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/coreos/go-semver/semver"

	"github.com/inconshreveable/go-update"
)

// baseURL is the source for auto-updates

var performAutoUpdate = false         // replaced at build time
var baseURL = "https://null.io/swarm" // replaced at build time

var internet InternetLike = &realInternet{}
var updater UpdaterLike = &realUpdater{}

// JSON represents the version.json file from the remote server
type JSON struct {
	Version string `json:"version"`
}

func versionURL() string {
	return baseURL + "/version.json"
}

// IsUpdateRequired determines whether an application update is needed
// by comparing the local version to the remote version
func IsUpdateRequired(localVersionString string) (bool, *semver.Version) {
	json := readStringFromURL(versionURL())
	remoteVersion := readJSONVersion(json)
	localVersion := semver.New(localVersionString)
	if localVersion.LessThan(*remoteVersion) {
		return true, remoteVersion
	}
	return false, nil
}

// AutoUpdate carries out an auto-update if needed and returns a bool indicating what occurred
func AutoUpdate(localVersionString string) (bool, *semver.Version) {
	if performAutoUpdate {
		updateReq, newVersion := IsUpdateRequired(localVersionString)
		if updateReq {
			fmt.Printf("Downloading version %v...", newVersion)
			downloadBytes, err := downloadBinary(newVersion)
			if err == nil {
				reader := bytes.NewReader(downloadBytes)

				err = updater.Apply(reader, update.Options{OldSavePath: oldSavePath()})
				if err == nil {
					return true, newVersion
				}
				log.Printf("Failed to apply update for version %v: %v", newVersion, err)
				return false, nil
			}
			log.Printf("Failed to download binary for %v: %v", newVersion, err)
			return false, nil
		}
	}

	// no update required
	return false, nil
}

func oldSavePath() string {
	oldSavePath := ".oldswarm"
	if runtime.GOOS == "windows" {
		oldSavePath += ".exe"
	}
	return oldSavePath
}

// downloadBinary triggers an update to the specified version
func downloadBinary(remoteVersion *semver.Version) ([]byte, error) {
	binaryURL := getBinaryURL(remoteVersion, runtime.GOOS)
	response, err := internet.Get(binaryURL)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	return bytes, err
}

func getBinaryURL(version *semver.Version, platform string) string {
	suffix := ".exe"
	if platform != "windows" {
		suffix = ""
	}
	url := fmt.Sprintf("%s/swarm-%s-%s%s", baseURL, version, platform, suffix)
	return url
}

// readJSONVersion returns the semver contained within a JSON string
func readJSONVersion(jsonString string) *semver.Version {
	var verJSON *JSON
	err := json.Unmarshal([]byte(jsonString), &verJSON)
	if err != nil {
		log.Printf("Error occurred reading JSON version: %v", err)
		return nil
	}
	version, err := semver.NewVersion(verJSON.Version)
	if err != nil {
		log.Printf("Error occurred reading JSON version: %v", err)
		return nil
	}
	return version
}

func readStringFromURL(url string) string {
	resp, err := internet.Get(url)
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	text := string(body)
	return text
}
