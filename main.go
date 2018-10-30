package main

import (
	"fmt"
	"swarm/bundle"
	"swarm/config"
	"swarm/monitor"
	"swarm/source"
	"swarm/ui"
	"swarm/util"
	"swarm/web"
)

func main() {
	ui.PrintTitle()

	// configuration
	swarmConfig, err := config.TryLoadSwarmConfigFromCWD()
	util.ExitIfError(err, "Failed to load swarm.json file: %s", err)
	runtimeConfig := ui.ChooseBuild(swarmConfig.Builds)
	moduleDescrs, err := config.LoadBuildDescriptionFile(runtimeConfig.BuildPath)
	util.ExitIfError(err, "Failed to load build description file: '%s'", runtimeConfig.BuildPath)

	// workspace
	ws := source.NewWorkspace(swarmConfig.RootPath)
	normalisedModules := moduleDescrs.NormaliseModules(ws.RootPath())
	moduleSet := bundle.CreateModuleSet(ws, normalisedModules, runtimeConfig)

	// web server
	handlers := moduleSet.GenerateHTTPHandlers()
	serverOptions := web.CreateServerOptions(swarmConfig.RootPath, swarmConfig.Server, handlers, runtimeConfig.BaseHref)
	server := web.CreateServer(serverOptions)
	hotReloader := web.NewHotReloader(server, ws, moduleSet)

	// monitor
	mon := monitor.NewMonitor(ws, swarmConfig.Monitor)
	mon.RegisterCallback(moduleSet.NotifyChanges)
	mon.RegisterCallback(hotReloader.NotifyReload)
	fmt.Print("Performing initial build...")
	mon.TriggerManually()

	go server.Start()
	go mon.NotifyOnChanges()
	fmt.Printf("Listening on %s\n", server.URL())
	if swarmConfig.Server.Open {
		util.OpenBrowser(server.URL())
	}

	// sleep
	util.WaitForCtrlC()
	server.Stop()
	mon.Stop()
}
