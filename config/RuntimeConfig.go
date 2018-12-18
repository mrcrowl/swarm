package config

// RuntimeConfig describes the expected state at runtime (currently, just what the base path will be)
type RuntimeConfig struct {
	// BaseHref gets the expected base path at runtime, e.g. <base href="app" /> ==> "app"
	BuildPath               string `json:"path"`
	BaseHref                string `json:"baseHref"`
	pathInterpolationValues map[string]string
}

// NewRuntimeConfig creates a RuntimeConfig
func NewRuntimeConfig(buildPath string, baseHref string) *RuntimeConfig {
	return &RuntimeConfig{buildPath, baseHref, map[string]string{}}
}

// SourceMapsEnabled ...
func (rtc *RuntimeConfig) SourceMapsEnabled() bool {
	return true
}

// SetPathInterpolationValues sets a map of key/value pairs to be interpolated into import paths
func (rtc *RuntimeConfig) SetPathInterpolationValues(values map[string]string) {
	rtc.pathInterpolationValues = values
}

// ImportPathInterpolationValues returns a map of key/value pairs to be interpolated into import paths
func (rtc *RuntimeConfig) ImportPathInterpolationValues() map[string]string {
	return rtc.pathInterpolationValues
}
