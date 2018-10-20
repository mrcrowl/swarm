package monitor

import (
	"gospm/dep"
	"log"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

// Monitor is used to recursively watch for file changes within a workspace
type Monitor struct {
	workspace *dep.Workspace
	channel   chan notify.EventInfo
}

// NewMonitor creates a new Monitor
func NewMonitor(workspace *dep.Workspace) *Monitor {
	channel := make(chan notify.EventInfo, 2048)

	rootRecursivePattern := filepath.Join(workspace.RootPath(), "./...")
	if err := notify.Watch(rootRecursivePattern, channel, notify.All); err != nil {
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
