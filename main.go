package main

import (
	"fmt"
	"log"
	"os"
	"swarm/bundle"
	"swarm/config"
	"swarm/monitor"
	"swarm/source"
	"swarm/ui"
	"swarm/util"
	"swarm/web"
)

func main() {
	log.SetOutput(os.Stdout)

	// title
	ui.PrintTitle()

	// configuration
	swarmConfig, err := config.TryLoadSwarmConfigFromCWD()
	exitIfError(err, "Failed to load swarm.json file: %s", err)

	runtimeConfig := ui.ChooseBuild(swarmConfig.Builds)
	moduleDescrs, err := config.LoadBuildDescriptionFile(runtimeConfig.BuildPath)
	exitIfError(err, "Failed to load build description file: '%s'", runtimeConfig.BuildPath)

	// workspace
	ws := source.NewWorkspace(swarmConfig.RootPath)
	normalisedModules := moduleDescrs.NormaliseModules(ws.RootPath())
	moduleSet := bundle.CreateModuleSet(ws, normalisedModules, runtimeConfig)

	// web server
	handlers := moduleSet.GenerateHTTPHandlers()
	server := web.CreateServer(swarmConfig.RootPath, &web.ServerOptions{
		Port:            swarmConfig.Server.Port,
		EnableHotReload: swarmConfig.Server.HotReload,
		Handlers:        handlers,
		BasePath:        runtimeConfig.BaseHref,
	})

	// monitor
	mon := monitor.NewMonitor(ws, swarmConfig.Monitor)
	mon.RegisterCallback(moduleSet.NotifyChanges)
	mon.RegisterCallback(server.NotifyReload)
	fmt.Print("Performing initial build...")
	mon.TriggerManually()

	go server.Start()
	go mon.NotifyOnChanges()
	fmt.Printf("Listening on %s", server.URL())
	if swarmConfig.Server.Open {
		util.OpenBrowser(server.URL())
	}

	// sleep
	util.WaitForExit()
	server.Stop()
}

func exitIfError(err error, message string, args ...interface{}) {
	if err != nil {
		log.Fatalf(message, args...)
		os.Exit(1)
	}
}
