package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"swarm/assets"
	"swarm/source"
	"swarm/util"
	"time"
)

const indexhtml = "index.html"
const systemJSConfigJS = "systemjs.config.js"
const swarmVirtualPath = "/__swarm__"
const assetsPhysicalPath = "/assets/static"
const hotReloadFilename = "HotReload.js"
const socketClientFilename = "SocketClient.js"
const webSocketServerPath = swarmVirtualPath + "/ws"

// Server is the state of the web server
type Server struct {
	srv          *http.Server
	rootFilepath string
	basePath     string
	port         uint16
	handlers     map[string]http.HandlerFunc
	hub          *SocketHub
}

// DefaultPort will be automatically assigned, if no port is specified in the options
const DefaultPort = uint16(8080)

// CreateServer returns a started webserver
func CreateServer(opts *ServerOptions) *Server {
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

	server := &Server{
		srv:          nil,
		rootFilepath: opts.RootFilepath,
		basePath:     opts.BasePath,
		port:         port,
		handlers:     opts.Handlers,
		hub:          hub,
	}

	return server
}

// Start the web server
func (server *Server) Start() {
	mux := http.NewServeMux()

	fileServer := server.attachStaticFileServer(mux)
	server.attachSystemJSRewriteHandler(mux)
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

func (server *Server) attachSystemJSRewriteHandler(mux *http.ServeMux) {
	systemJSFilepath := filepath.Join(server.rootFilepath, server.basePath, systemJSConfigJS)
	handler := func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile(systemJSFilepath)
		if err != nil {
			log.Printf("ERROR: Failed to load SystemJS config at: %s", systemJSFilepath)
			return
		}
		configJS := string(bytes)
		rewrittenConfigJS := rewriteSystemJSConfigPaths(configJS)
		mimeType := util.MimeTypeFromFilename(systemJSFilepath)
		w.Header().Set("Content-Type", mimeType)
		io.WriteString(w, rewrittenConfigJS)
		return
	}
	systemJSPath := path.Join("/", server.basePath, systemJSConfigJS)
	mux.HandleFunc(systemJSPath, handler)
}

var rewriteSystemJSPattern = regexp.MustCompile(`"\.\/(common|services|utils)",`)

func rewriteSystemJSConfigPaths(systemJSConfig string) string {
	return rewriteSystemJSPattern.ReplaceAllString(systemJSConfig, `"../$1", /* <-- REWRITTEN BY SWARM */`)
}

func loadAssetString(assetFilename string) string {
	source, _ := assets.FS.String(assetsPhysicalPath + "/" + assetFilename)
	return source
}

func swarmify(assetPath string) string {
	return path.Join(swarmVirtualPath, assetPath)
}

func createStringHandleFunc(filename string) func(w http.ResponseWriter, r *http.Request) {
	contents := loadAssetString(filename)
	mimeType := util.MimeTypeFromFilename(filename)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", mimeType)
		io.WriteString(w, contents)
	}
}

func (server *Server) attachWebSocketListeners(mux *http.ServeMux, hub *SocketHub) {
	mux.HandleFunc(swarmify(socketClientFilename), createStringHandleFunc(socketClientFilename))
	mux.HandleFunc(swarmify(hotReloadFilename), createStringHandleFunc(hotReloadFilename))
	mux.HandleFunc(webSocketServerPath, func(w http.ResponseWriter, r *http.Request) {
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
				indexHTML := string(bytes)
				injectedIndexHTML := InjectSrcJavascript(indexHTML, swarmify(hotReloadFilename), true)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.WriteString(w, injectedIndexHTML)
				return
			}
		}

		fileServer.ServeHTTP(w, r)
	}

	mux.HandleFunc(rootedBasePath+"/", indexHandler)
}

// TriggerFullReload causes a full HTML reload to be fired
func (server *Server) TriggerFullReload() {
	server.hub.broadcast("reload", "")
}

// ReloadCSSPayloadData encapsulates the data to reload a specific style sheet
type ReloadCSSPayloadData struct {
	ID  string `json:"id"`
	CSS string `json:"css"`
}

// TriggerCSSReload causes a CSS-only reload to be fired
func (server *Server) TriggerCSSReload(path string, css string) {
	cssReloadData := &ReloadCSSPayloadData{
		ID:  source.CSSPrefix + path,
		CSS: css,
	}
	jsonBytes, _ := json.Marshal(cssReloadData)
	server.hub.broadcast("reload-css", string(jsonBytes))
}

// URL gets the localhost URL for this server
func (server *Server) URL() string {
	return fmt.Sprintf("http://localhost:%d/%s", server.Port(), server.basePath)
}

// IsHotReloadEnabled gets whether hot reload is enabled
func (server *Server) IsHotReloadEnabled() bool {
	return server.hub != nil
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
	}
	if server.hub != nil {
		server.hub.stop()
	}
}
