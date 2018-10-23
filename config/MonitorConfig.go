package config

// MonitorConfig describes the configuration of the file monitor
type MonitorConfig struct {
	Extensions     []string `json:"extensions"`
	DebounceMillis uint     `json:"debounceMillis"`
}

// NewMonitorConfig creates a MonitorConfig
func NewMonitorConfig(extensions []string, debounceMillis uint) *MonitorConfig {
	return &MonitorConfig{extensions, debounceMillis}
}
