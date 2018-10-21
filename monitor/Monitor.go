package monitor

import (
	"gospm/source"
	"log"
	"path/filepath"
	"time"

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

const notifyInterval = 10 * time.Minute
const debounceInterval = 500 * time.Millisecond

// NotifyOnChanges notifies when events occur (after debouncing)
func (mon *Monitor) NotifyOnChanges(callback func(changes *EventChangeset)) {
	debounceTimer := time.NewTimer(notifyInterval)
	changeset := NewEventChangeset()

	var e notify.EventInfo
	for {
		select {
		case e = <-mon.channel:
			// receive an event
			changeset.Add(e.Event(), e.Path())
			debounceTimer.Reset(debounceInterval)

		case <-debounceTimer.C:
			// debounce and fire callback
			if changeset.nonEmpty() {
				callback(changeset)
			}
			changeset = NewEventChangeset()
		}
	}
}

// Stop cancels the recursive watcher
func (mon *Monitor) Stop() {
	notify.Stop(mon.channel)
}
