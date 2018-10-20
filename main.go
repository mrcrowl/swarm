package main

import (
	"gospm/bundle"
	"gospm/dep"
	"gospm/source"
	"log"
)

func main() {
	ws := source.NewWorkspace("c:\\wf\\lp\\web\\App")
	fileset := dep.BuildFileSet(ws, "app/src/ep/app")

	log.Println(fileset.Count())

	bundler := bundle.NewBundler()
	bundler.Bundle(fileset)
}

func monitor() {
	// wsw := dep.NewWorkspaceWatcher(ws)

	// for {
	// 	eventInfo := <-wsw.C()
	// 	log.Println("Got event: ", eventInfo)
	// }
}
