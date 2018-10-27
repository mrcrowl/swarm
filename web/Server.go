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
	hub      *SocketHub
}

// ServerOptions specifies the parameters for the web server
type ServerOptions struct {
	Port     uint16
	Handlers map[string]http.HandlerFunc
}

var serverInstance *Server

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// CreateServer returns a started webserver
func CreateServer(rootPath string, opts *ServerOptions) *Server {
	if serverInstance != nil {
		return serverInstance
	}

	port := DefaultPort
	if opts != nil {
		port = opts.Port
	}

	hub := newSocketHub()

	serverInstance = &Server{
		srv:      nil,
		rootPath: rootPath,
		port:     port,
		handlers: opts.Handlers,
		hub:      hub,
	}

	return serverInstance
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

	// WEBSOCKETS
	mux.HandleFunc("/__swarm__/SocketClient.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./web/static/SocketClient.js")
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(server.hub, w, r)
	})

	go server.hub.run()
	// ENDWEBSOCKETS

	server.srv = &http.Server{
		Addr:    makeServerAddress(server.port),
		Handler: mux,
	}

	if err := server.srv.ListenAndServe(); err != nil {
		panic(err)
	}

}

// NotifyReload sends a message to the client page to reload
func (server *Server) NotifyReload() {
	server.hub.broadcast("reload", "")
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
		serverInstance = nil
	}
}
