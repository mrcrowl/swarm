package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const swarmConfigDefaultFilename = "swarm.json"
const defaultRootPathWindows = "c:\\wf\\lp\\web\\App"
const defaultRootPathMacOSAndLinux = "%s/web/App" // <-- %s will be replaced with user dir
const defaultServerPort uint16 = 8096

var defaultMonitorExtensions = []string{".js", ".html", ".css", ".json"}
var defaultBuilds = map[string]*RuntimeConfig{
	"app": NewRuntimeConfig(
		"build/systemjs_build_app.json",
		"app",
	),
	"controlpanel": NewRuntimeConfig(
		"build/systemjs_build_controlpanel.json",
		"controlpanel",
	),
}

// SwarmConfig is the root configuration file
type SwarmConfig struct {
	RootPath string                    `json:"root"`
	Monitor  *MonitorConfig            `json:"monitor"`
	Builds   map[string]*RuntimeConfig `json:"builds"`
	Server   *ServerConfig             `json:"server"`
}

func (config *SwarmConfig) expandAndNormalisePaths(cwd string) {
	norm := func(base string, path string) string {
		if filepath.IsAbs(path) {
			return filepath.Clean(path)
		}

		return filepath.Join(base, path)
	}

	config.RootPath = norm(cwd, config.RootPath)
	for _, b := range config.Builds {
		b.BuildPath = norm(config.RootPath, b.BuildPath)
	}
}

func (config *SwarmConfig) backfillWithDefaults(cwd string) {
	defaults := DefaultSwarmConfig(cwd)

	if config.RootPath == "" {
		config.RootPath = defaults.RootPath
	}

	if config.Builds == nil {
		config.Builds = defaults.Builds
	}

	if config.Monitor == nil {
		config.Monitor = defaults.Monitor
	}

	if config.Server == nil {
		config.Server = defaults.Server
	}
}

// TryLoadSwarmConfigFromCWD tries to load a swarm.json configuration from the current working directory
func TryLoadSwarmConfigFromCWD() (*SwarmConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get CWD: %s", err)
	}
	swarmjsonFilepath := filepath.Join(cwd, swarmConfigDefaultFilename)
	if _, err := os.Stat(swarmjsonFilepath); err != nil {
		return DefaultSwarmConfig(cwd), nil
	}
	return LoadSwarmConfig(swarmjsonFilepath, cwd)
}

func getDefaultRootPath() string {
	if runtime.GOOS == "windows" {
		return defaultRootPathWindows
	}

	usr, _ := user.Current()
	return fmt.Sprintf(defaultRootPathMacOSAndLinux, usr.HomeDir)
}

// DefaultSwarmConfig loads the default swarm.json configuration file
func DefaultSwarmConfig(cwd string) *SwarmConfig {
	config := &SwarmConfig{
		RootPath: getDefaultRootPath(),
		Monitor:  NewMonitorConfig(defaultMonitorExtensions, 150),
		Builds:   defaultBuilds,
		Server:   NewServerConfig(defaultServerPort, true),
	}
	config.expandAndNormalisePaths(cwd)
	return config
}

// LoadSwarmConfig loads a swarm.json configuration file
func LoadSwarmConfig(swarmConfigFilepath string, cwd string) (*SwarmConfig, error) {
	if filepath.Ext(swarmConfigFilepath) == "" {
		swarmConfigFilepath += ".json"
	}

	buildBytes, e := ioutil.ReadFile(swarmConfigFilepath)
	if e != nil {
		return nil, errors.New("Invalid swarm config file: " + swarmConfigFilepath)
	}

	jsonString := string(buildBytes)
	config, err := LoadSwarmConfigString(jsonString, cwd)
	return config, err
}

// LoadSwarmConfigString loads a swarm.json file from a string
func LoadSwarmConfigString(swarmConfigJSON string, cwd string) (*SwarmConfig, error) {
	var config *SwarmConfig
	err := json.Unmarshal([]byte(swarmConfigJSON), &config)
	if err != nil {
		return nil, errors.New("Invalid JSON in swarm config file: " + err.Error())
	}
	config.backfillWithDefaults(cwd)
	config.expandAndNormalisePaths(cwd)
	return config, nil
}
