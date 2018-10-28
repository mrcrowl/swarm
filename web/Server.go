package web

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"swarm/monitor"
	"time"
)

const indexhtml = "index.html"
const reloadJavascript = `
		import {SocketClient} from "/__swarm__/SocketClient.js"
		const sc = new SocketClient()
		sc.on(e => e.type == "reload" && window.location.reload());
		sc.connect();
`
const socketClientPath = "/__swarm__/SocketClient.js"
const socketClientSource = "./web/static/SocketClient.js"

// Server is the state of the web server
type Server struct {
	srv          *http.Server
	rootFilepath string
	basePath     string
	port         uint16
	handlers     map[string]http.HandlerFunc
	hub          *SocketHub
}

var serverInstance *Server

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// CreateServer returns a started webserver
func CreateServer(opts *ServerOptions) *Server {
	if serverInstance != nil {
		return serverInstance
	}

	port := DefaultPort
	if opts != nil {
		port = opts.Port
	}

	enableHotReload := true
	if opts != nil {
		enableHotReload = opts.EnableHotReload
	}

	hub := (*SocketHub)(nil)
	if enableHotReload {
		hub = newSocketHub()
	}

	serverInstance = &Server{
		srv:          nil,
		rootFilepath: opts.RootFilepath,
		basePath:     opts.BasePath,
		port:         port,
		handlers:     opts.Handlers,
		hub:          hub,
	}

	return serverInstance
}

// Start the web server
func (server *Server) Start() {
	mux := http.NewServeMux()

	fileServer := server.attachStaticFileServer(mux)
	server.attachCustomHandlers(mux)

	if server.hub != nil {
		// add HMR support
		server.attachIndexInjectionListener(mux, fileServer)
		server.attachWebSocketListeners(mux, server.hub)
		go server.hub.run()
	}

	server.srv = &http.Server{
		Addr:    makeServerAddress(server.port),
		Handler: mux,
	}

	if err := server.srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (server *Server) attachCustomHandlers(mux *http.ServeMux) {
	for url, handler := range server.handlers {
		mux.HandleFunc(url, handler)
	}
}

func (server *Server) attachStaticFileServer(mux *http.ServeMux) http.Handler {
	fileServer := http.FileServer(http.Dir(server.rootFilepath))
	mux.Handle("/", fileServer)
	return fileServer
}

func (server *Server) attachWebSocketListeners(mux *http.ServeMux, hub *SocketHub) {
	mux.HandleFunc(socketClientPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		http.ServeFile(w, r, socketClientSource)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(hub, w, r)
	})
}

func (server *Server) attachIndexInjectionListener(mux *http.ServeMux, fileServer http.Handler) {
	rootedBasePath := path.Join("/", server.basePath)
	acceptedIndexPaths := []string{
		rootedBasePath,
		rootedBasePath + "/",
		rootedBasePath + "/" + indexhtml,
	}

	indexFilepath := filepath.Join(server.rootFilepath, server.basePath, indexhtml)
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		for _, path := range acceptedIndexPaths {
			if r.URL.Path == path {
				bytes, err := ioutil.ReadFile(indexFilepath)
				if err != nil {
					log.Printf("ERROR: Failed to load index at: %s", indexFilepath)
					return
				}
				indexString := string(bytes)
				injectedIndexString := InjectInlineJavascript(indexString, reloadJavascript, true)
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, injectedIndexString)
				return
			}
		}

		fileServer.ServeHTTP(w, r)
	}

	mux.HandleFunc(rootedBasePath+"/", indexHandler)
}

// NotifyReload sends a message to the client page to reload
func (server *Server) NotifyReload(changes *monitor.EventChangeset) {
	if server.hub != nil {
		server.hub.broadcast("reload", "")
	}
}

// URL gets the localhost URL for this server
func (server *Server) URL() string {
	return fmt.Sprintf("http://localhost:%d/%s", server.Port(), server.basePath)
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
