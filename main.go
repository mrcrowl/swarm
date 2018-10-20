package main

import (
	"gospm/dep"
	"log"
)

func main() {
	ws := dep.NewWorkspace("c:\\wf\\lp\\web\\App")
	wsw := dep.NewWorkspaceWatcher(ws)

	for {
		eventInfo := <-wsw.C()
		log.Println("Got event: ", eventInfo)
	}
}
