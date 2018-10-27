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
	"time"
)

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
	indexPath    string
	port         uint16
	handlers     map[string]http.HandlerFunc
	hub          *SocketHub
}

// ServerOptions specifies the parameters for the web server
type ServerOptions struct {
	Port      uint16
	Handlers  map[string]http.HandlerFunc
	IndexPath string
}

var serverInstance *Server

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// CreateServer returns a started webserver
func CreateServer(rootFilepath string, opts *ServerOptions) *Server {
	if serverInstance != nil {
		return serverInstance
	}

	port := DefaultPort
	if opts != nil {
		port = opts.Port
	}

	hub := newSocketHub()

	serverInstance = &Server{
		srv:          nil,
		rootFilepath: rootFilepath,
		indexPath:    opts.IndexPath,
		port:         port,
		handlers:     opts.Handlers,
		hub:          hub,
	}

	return serverInstance
}

// Start the web server
func (server *Server) Start() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(server.rootFilepath))
	mux.Handle("/", fileServer)
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

	server.attachIndexInjectionListener(mux, fileServer)
	server.attachWebSocketListeners(mux)
	go server.hub.run()

	server.srv = &http.Server{
		Addr:    makeServerAddress(server.port),
		Handler: mux,
	}

	if err := server.srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (server *Server) attachWebSocketListeners(mux *http.ServeMux) {
	mux.HandleFunc(socketClientPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		http.ServeFile(w, r, socketClientSource)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(server.hub, w, r)
	})
}

func (server *Server) attachIndexInjectionListener(mux *http.ServeMux, fileServer http.Handler) {
	indexDirPath := path.Join("/", path.Dir(server.indexPath))
	acceptedIndexPaths := []string{
		"/" + server.indexPath,
		indexDirPath,
		indexDirPath + "/",
	}

	indexFilepath := filepath.Join(server.rootFilepath, server.indexPath)
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

	mux.HandleFunc(indexDirPath+"/", indexHandler)
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
