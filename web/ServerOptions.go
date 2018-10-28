package web

import (
	"net/http"
	"swarm/config"
)

// ServerOptions specifies the parameters for the web server
type ServerOptions struct {
	RootFilepath    string
	Port            uint16
	EnableHotReload bool
	Handlers        map[string]http.HandlerFunc
	BasePath        string
}

// CreateServerOptions forms a server options object from various sources
func CreateServerOptions(
	rootFilepath string,
	serverConfig *config.ServerConfig,
	handlers map[string]http.HandlerFunc,
	basePath string,
) *ServerOptions {
	return &ServerOptions{
		RootFilepath:    rootFilepath,
		Port:            serverConfig.Port,
		EnableHotReload: serverConfig.HotReload,
		Handlers:        handlers,
		BasePath:        basePath,
	}
}
