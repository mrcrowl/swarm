package main

import (
	"gospm/bundle"
	"gospm/dep"
	"gospm/monitor"
	"gospm/source"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ws := source.NewWorkspace("c:\\wf\\lp\\web\\App")
	// ws := source.NewWorkspace("c:\\wf\\home\\topo\\")
	// fileset := source.NewEmptyFileSet()
	fileset := dep.BuildFileSet(ws, "app/src/ep/app")

	mon := monitor.NewMonitor(ws)

	makeBundle := func(_ *monitor.EventChangeset) {
		log.Print("Bundling...")
		bundler := bundle.NewBundler()
		bundler.Bundle(fileset)
		log.Println("Done")
	}

	go mon.NotifyOnChanges(makeBundle)
	makeBundle(nil)

	// waitForExit()
}

func waitForExit() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	// systemTeardown()
}
