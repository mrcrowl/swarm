package version

import (
	"runtime"
	"testing"

	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/assert"
)

func TestReadJSONVersion(t *testing.T) {
	actual := readJSONVersion(`{"version": "1.2.3"}`)
	assert.NotNil(t, actual)
	expected := semver.New("1.2.3")
	assert.True(t, expected.Equal(*actual), "Expected %v to be %v", actual, expected)
}

func TestReadJSONVersionInvalidJSON(t *testing.T) {
	actual := readJSONVersion(`{"version": 1.2.3"}`)
	assert.Nil(t, actual)
}

func TestReadJSONVersionInvalidVersion(t *testing.T) {
	actual := readJSONVersion(`{"version": ""}`)
	assert.Nil(t, actual)
}

func TestReadStringFromURL(t *testing.T) {
	contents := readStringFromURL(versionURL)
	assert.NotEmpty(t, contents)
}

func TestReadStringFromURLBadURL(t *testing.T) {
	contents := readStringFromURL("no.protocol.specified.com/swarm/version.json")
	assert.Empty(t, contents)
}

func TestIsUpdateRequiredTrue(t *testing.T) {
	mock := newMockInternet()
	mock.addStringResponse(versionURL, `{ "version": "1.0.0" }`)
	internet = mock
	actual, version := IsUpdateRequired("0.9.1")
	assert.NotNil(t, version)
	assert.Equal(t, true, actual)
}

func TestIsUpdateRequiredFalse(t *testing.T) {
	mock := newMockInternet()
	mock.addStringResponse(versionURL, `{ "version": "1.0.0" }`)
	internet = mock
	actual, version := IsUpdateRequired("1.0.0")
	assert.Nil(t, version)
	assert.Equal(t, false, actual)
}

func TestDownloadBinary(t *testing.T) {
	remoteVersion := semver.New("1.0.1")
	platform := runtime.GOOS
	binaryURL := getBinaryURL(remoteVersion, platform)
	mock := newMockInternet()
	mock.addStringResponse(binaryURL, string([]byte{0x1, 0x2, 0x3, 0x4}))
	mock.addStringResponse(versionURL, `{ "version": "1.0.1" }`)
	internet = mock

	bytes, err := downloadBinary(remoteVersion)
	assert.Nil(t, err)
	assert.Len(t, bytes, 4)
}

func TestGetBinaryURL(t *testing.T) {
	cases := map[string]struct {
		version  string
		platform string
		expected string
	}{
		"window": {
			version:  "1.9.0",
			platform: "windows",
			expected: versionsPath + "/swarm-1.9.0-windows.exe",
		},
		"macOS": {
			version:  "1.1.0",
			platform: "darwin",
			expected: versionsPath + "/swarm-1.1.0-darwin",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			version := semver.New(tc.version)
			actual := getBinaryURL(version, tc.platform)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestAutoUpdate(t *testing.T) {
	remoteVersion := semver.New("1.0.1")
	binaryURL := getBinaryURL(remoteVersion, runtime.GOOS)
	mock := newMockInternet()
	mock.addStringResponse(binaryURL, string([]byte{0x1, 0x2, 0x3, 0x4}))
	mock.addStringResponse(versionURL, `{ "version": "1.0.1" }`)
	internet = mock

	mockUpdater := &mockUpdater{}
	updater = mockUpdater

	cases := map[string]struct {
		localVersion string
		didUpdate    bool
	}{
		"lower":  {localVersion: "1.0.0", didUpdate: true},
		"equal":  {localVersion: "1.0.1", didUpdate: false},
		"higher": {localVersion: "1.0.2", didUpdate: false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			localVersion := tc.localVersion

			updated, newVersion := AutoUpdate(localVersion)
			assert.Equal(t, tc.didUpdate, updated)
			if updated {
				assert.True(t, remoteVersion.Equal(*newVersion), "Expected %v to be %v", newVersion, remoteVersion)
				assert.True(t, mockUpdater.successful)
			}
		})
	}
}
