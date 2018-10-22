package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"swarm/bundle"
	"swarm/dep"
	"swarm/monitor"
	"swarm/source"
	"swarm/web"
	"sync"
	"syscall"

	"github.com/rjeczalik/notify"
)

const folder = "c:\\wf\\lp\\web\\App"

const app = folder + "\\app\\src\\ep\\app.js"

func main() {
	log.SetOutput(os.Stdout)

	ws := source.NewWorkspace(folder)
	filterFn := func(event notify.Event, path string) bool {
		if strings.HasSuffix(path, ".ts") {
			return false
		}
		return true
	}
	mon := monitor.NewMonitor(ws, filterFn)

	var appjs string
	bundleMutex := &sync.Mutex{}
	makeBundle := func(changeset *monitor.EventChangeset) {
		fileset := dep.BuildFileSet(ws, "app/src/ep/app")
		bundleMutex.Lock()
		defer bundleMutex.Unlock()

		log.Print("Bundling...")
		bundler := bundle.NewBundler()
		sb := bundler.Bundle(fileset)
		appjs = sb.String()
		// ioutil.WriteFile(app, []byte(sb.String()), os.ModePerm) // HACK
		log.Println("Done")
	}

	go mon.NotifyOnChanges(makeBundle)
	makeBundle(nil)

	handlers := map[string]http.HandlerFunc{
		"/app/src/ep/app.js": func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, appjs)
		},
	}

	server := web.CreateServer(folder, &web.ServerOptions{
		Port:     8096,
		Handlers: handlers,
	})
	go server.Start()
	log.Printf("Listening on http://localhost:%d", server.Port())
	waitForExit(server)
}

func waitForExit(server *web.Server) {

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	server.Stop()
}
