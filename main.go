package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"swarm/bundle"
	"swarm/config"
	"swarm/monitor"
	"swarm/source"
	"swarm/web"
	"syscall"
)

const folder = "c:\\wf\\lp\\web\\App"
const app = folder + "\\app\\src\\ep\\app.js"
const buildFile = "C:\\WF\\LP\\web\\App\\build\\systemjs_build_app.json"

func main() {
	log.SetOutput(os.Stdout)

	swarmConfig, err := config.TryLoadSwarmConfigFromCWD()
	if err != nil {
		log.Fatalf("Failed to load swarm.json file: %s", err)
		return
	}

	ws := source.NewWorkspace(folder)
	mon := monitor.NewMonitor(ws, swarmConfig.Monitor)

	moduleDescrs, err := config.LoadBuildDescriptionFile(buildFile)
	if err != nil {
		log.Fatalf("Failed to load build description file: '%s'", buildFile)
		return
	}

	selectedBuild := "app"
	runtimeConfig := swarmConfig.Builds[selectedBuild]

	moduleSet := bundle.CreateModuleSet(
		ws,
		moduleDescrs.NormaliseModules(ws.RootPath()),
		runtimeConfig,
	)

	go mon.NotifyOnChanges(moduleSet.NotifyChanges)
	moduleSet.NotifyChanges(nil)

	handlers := moduleSet.GenerateHTTPHandlers()

	server := web.CreateServer(folder, &web.ServerOptions{
		Port:     swarmConfig.Server.Port,
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
