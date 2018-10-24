package config

// RuntimeConfig describes the expected state at runtime (currently, just what the base path will be)
type RuntimeConfig struct {
	// BaseHref gets the expected base path at runtime, e.g. <base href="app" /> ==> "app"
	BuildPath string `json:"path"`
	BaseHref  string `json:"baseHref"`
}

// NewRuntimeConfig creates a RuntimeConfig
func NewRuntimeConfig(buildPath string, baseHref string) *RuntimeConfig {
	return &RuntimeConfig{buildPath, baseHref}
}

// SourceMapsEnabled ...
func (rtc *RuntimeConfig) SourceMapsEnabled() bool {
	return true // TEMPORARILY DISABLED UNTIL I CAN GET IT PERFORMING NICELY :(
}
