package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"swarm/bundle"
	"swarm/dep"
	"swarm/monitor"
	"swarm/source"
	"syscall"
)

const folder = "c:\\wf\\lp\\web\\App"

// const folder = "c:\\wf\\home\\topo"
const app = folder + "\\app\\src\\ep\\app.js"
const appmoved = folder + "\\app\\src\\ep\\app-moved.js"

func main() {
	log.SetOutput(os.Stdout)

	// HACK
	if _, err := os.Stat(appmoved); err == nil {
		os.Remove(app)
		os.Rename(appmoved, app)
	}
	// ENDHACK

	ws := source.NewWorkspace(folder)
	// ws := source.NewWorkspace("c:\\wf\\home\\topo\\")
	// fileset := source.NewEmptyFileSet()
	fileset := dep.BuildFileSet(ws, "app/src/ep/app")

	mon := monitor.NewMonitor(ws)

	makeBundle := func(_ *monitor.EventChangeset) {
		log.Print("Bundling...")
		bundler := bundle.NewBundler()
		sb := bundler.Bundle(fileset)
		os.Rename(app, appmoved)
		ioutil.WriteFile(app, []byte(sb.String()), os.ModePerm) // HACK
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
