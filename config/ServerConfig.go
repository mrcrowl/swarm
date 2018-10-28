package config

// ServerConfig is the configuration for the built-in web server
type ServerConfig struct {
	Port      uint16
	HotReload bool
}

// NewServerConfig creates a new ServerConfig
func NewServerConfig(port uint16, enableHotReload bool) *ServerConfig {
	return &ServerConfig{port, enableHotReload}
}
