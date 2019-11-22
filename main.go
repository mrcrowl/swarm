package main

import (
	"fmt"
	"os"

	"github.com/mrcrowl/swarm/bundle"
	"github.com/mrcrowl/swarm/config"
	"github.com/mrcrowl/swarm/monitor"
	"github.com/mrcrowl/swarm/source"
	"github.com/mrcrowl/swarm/ui"
	"github.com/mrcrowl/swarm/util"
	"github.com/mrcrowl/swarm/version"
	"github.com/mrcrowl/swarm/web"

	flag "github.com/spf13/pflag"
)

const localver = "1.0.11"

var portFlag = flag.Uint16P("port", "p", uint16(8096), "Web server port number")
var helpFlag = flag.BoolP("help", "h", false, "Shows the usage")

func main() {
	ui.PrintTitle(localver)
	ui.CheckHelp(helpFlag)

	if didUpdate, _ := version.AutoUpdate(localver); didUpdate {
		fmt.Println("updated. Please restart!")
		os.Exit(0)
	}

	// configuration
	swarmConfig, err := config.TryLoadSwarmConfigFromCWD(portFlag)
	util.ExitIfError(err, "Failed to load swarm.json file: %s", err)
	runtimeConfig := ui.ChooseBuild(swarmConfig.Builds, flag.Args())
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
