package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"swarm/bundle"
	"swarm/config"
	"swarm/monitor"
	"swarm/source"
	"swarm/web"
	"syscall"

	"github.com/rjeczalik/notify"
)

const folder = "c:\\wf\\lp\\web\\App"
const app = folder + "\\app\\src\\ep\\app.js"
const buildFile = "C:\\WF\\LP\\web\\App\\build\\systemjs_build_app.json"

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

	moduleDescrs, err := config.LoadBuildDescriptionFile(buildFile)
	if err != nil {
		log.Fatalf("Failed to load build description file: '%s'", buildFile)
	}

	moduleSet := bundle.CreateModuleSet(ws, moduleDescrs.NormaliseModules(ws.RootPath()))

	go mon.NotifyOnChanges(moduleSet.NotifyChanges)
	moduleSet.NotifyChanges(nil)

	handlers := moduleSet.GenerateHTTPHandlers()

	server := web.CreateServer(folder, &web.ServerOptions{
		Port:     8096,
		Handlers: handlers,
	})
	go server.Start()
	fmt.Printf("Listening on http://localhost:%d\n", server.Port())
	waitForExit(server)
}

func waitForExit(server *web.Server) {

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	server.Stop()
}
