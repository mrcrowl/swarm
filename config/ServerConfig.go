package config

// ServerConfig is the configuration for the built-in web server
type ServerConfig struct {
	Port      uint16 `json:"port"`
	Open      bool   `json:"open"`
	HotReload bool   `json:"hotReload"`
}

// NewServerConfig creates a new ServerConfig
func NewServerConfig(port uint16, open bool, enableHotReload bool) *ServerConfig {
	return &ServerConfig{port, open, enableHotReload}
}
