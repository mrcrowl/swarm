package monitor

import (
	"gospm/source"
	"log"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

// Monitor is used to recursively watch for file changes within a workspace
type Monitor struct {
	workspace *source.Workspace
	channel   chan notify.EventInfo
}

// NewMonitor creates a new Monitor
func NewMonitor(workspace *source.Workspace) *Monitor {
	channel := make(chan notify.EventInfo, 2048)

	rootPathRecursive := filepath.Join(workspace.RootPath(), "./...")
	if err := notify.Watch(rootPathRecursive, channel, (notify.Write | notify.Remove)); err != nil {
		log.Fatal(err)
	}

	return &Monitor{
		workspace,
		channel,
	}
}

// C returns the channel which notifies about events
func (mon *Monitor) C() chan notify.EventInfo {
	return mon.channel
}

// Stop cancels the recursive watcher
func (mon *Monitor) Stop() {
	notify.Stop(mon.channel)
}
