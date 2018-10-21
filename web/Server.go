package web

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server is the state of the web server
type Server struct {
	server   *http.Server
	rootPath string
	port     uint16
}

// ServerOptions specifies the parameters for the web server
type ServerOptions struct {
	Port uint16
}

var stateInstance *Server

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// StartServer returns a started webserver
func StartServer(rootPath string, opts *ServerOptions) *Server {
	if stateInstance != nil {
		return stateInstance
	}

	port := DefaultPort
	if opts != nil {
		port = opts.Port
	}

	stateInstance = &Server{
		server:   nil,
		rootPath: rootPath,
		port:     port,
	}

	stateInstance.start()
	return stateInstance
}

// start the web server
func (state *Server) start() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(state.rootPath)))

	state.server = &http.Server{
		Addr:    makeServerAddress(state.port),
		Handler: mux,
	}

	if err := state.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func makeServerAddress(port uint16) string {
	return fmt.Sprintf(":%d", port)
}

// Stop the web server
func (state *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if cancel != nil {
		state.server.Shutdown(ctx)
		state.server = nil
		stateInstance = nil
	}
}
