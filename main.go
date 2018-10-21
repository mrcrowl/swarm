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
	fileset := dep.BuildFileSet(ws, "app/src/ep/app")

	mon := monitor.NewMonitor(ws)
	go mon.NotifyOnChanges(func(_ *monitor.EventChangeset) {
		log.Println("Bundling...")
		bundler := bundle.NewBundler()
		bundler.Bundle(fileset)
	})

	waitForExit()
}

func waitForExit() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	// systemTeardown()
}
