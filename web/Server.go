package web

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server is the state of the web server
type Server struct {
	srv      *http.Server
	rootPath string
	port     uint16
	handlers map[string]http.HandlerFunc
}

// ServerOptions specifies the parameters for the web server
type ServerOptions struct {
	Port     uint16
	Handlers map[string]http.HandlerFunc
}

var stateInstance *Server

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// CreateServer returns a started webserver
func CreateServer(rootPath string, opts *ServerOptions) *Server {
	if stateInstance != nil {
		return stateInstance
	}

	port := DefaultPort
	if opts != nil {
		port = opts.Port
	}

	stateInstance = &Server{
		srv:      nil,
		rootPath: rootPath,
		port:     port,
		handlers: opts.Handlers,
	}

	return stateInstance
}

// Start the web server
func (server *Server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(server.rootPath)))
	for url, handler := range server.handlers {
		mux.HandleFunc(url, handler)
	}

	// HACK:
	// mux.HandleFunc("/app/src/ep/app.js.map", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "C:\\WF\\Home\\topo\\app\\src\\ep\\app.build.js.map")
	// })
	// mux.HandleFunc("/app/src/ep/ui/base/BaseController.ts", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "C:\\WF\\LP\\Web\\App\\app\\src\\ep\\ui\\base\\BaseController.ts")
	// })
	// ENDHACK

	server.srv = &http.Server{
		Addr:    makeServerAddress(server.port),
		Handler: mux,
	}

	if err := server.srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Port gets the port number for this server
func (server *Server) Port() uint16 {
	return server.port
}

func makeServerAddress(port uint16) string {
	return fmt.Sprintf(":%d", port)
}

// Stop the web server
func (server *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if cancel != nil {
		server.srv.Shutdown(ctx)
		server.srv = nil
		stateInstance = nil
	}
}
