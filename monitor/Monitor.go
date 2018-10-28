package monitor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"swarm/config"
	"swarm/source"
	"sync"
	"time"

	"github.com/rjeczalik/notify"
)

// FilterFn is the shape of a function that can used as a filter for a Monitor
type FilterFn func(notify.Event, string) bool

// Monitor is used to recursively watch for file changes within a workspace
type Monitor struct {
	workspace        *source.Workspace
	channel          chan notify.EventInfo
	filter           FilterFn
	debounceDuration time.Duration
	changeCallbacks  []func(changes *EventChangeset)
	callbackMutex    *sync.Mutex
}

// NewMonitor creates a new Monitor
func NewMonitor(workspace *source.Workspace, config *config.MonitorConfig) *Monitor {
	channel := make(chan notify.EventInfo, 2048)

	rootPathRecursive := filepath.Join(workspace.RootPath(), "./...")
	if err := notify.Watch(rootPathRecursive, channel, (notify.Write | notify.Remove)); err != nil {
		log.Fatal(err)
	}

	filter := createExtensionFilterFn(config.Extensions)
	debounceDuration := time.Millisecond * time.Duration(config.DebounceMillis)
	callbackMutex := &sync.Mutex{}

	return &Monitor{
		workspace,
		channel,
		filter,
		debounceDuration,
		nil,
		callbackMutex,
	}
}

func createExtensionFilterFn(extensions []string) FilterFn {
	return func(event notify.Event, path string) bool {
		ext := filepath.Ext(path)
		for _, validExt := range extensions {
			if ext == validExt {
				return true
			}
		}

		return false
	}
}

const notifyInterval = 10 * time.Minute

// RegisterCallback adds a callback function which will be called when a change occurs
func (mon *Monitor) RegisterCallback(callback func(changes *EventChangeset)) {
	mon.changeCallbacks = append(mon.changeCallbacks, callback)
}

// TriggerManually is used to manually trigger the NotifyOnChanges event
func (mon *Monitor) TriggerManually() {
	mon.triggerCallbacks(nil, time.Now(), true)
}

func (mon *Monitor) triggerCallbacks(changeset *EventChangeset, start time.Time, silent bool) {
	mon.callbackMutex.Lock()

	fmt.Println("")
	for _, callback := range mon.changeCallbacks {
		callback(changeset)
	}

	if !silent {
		elapsed := time.Since(start)
		fmt.Printf("...done in %s\n", elapsed)
	}

	mon.callbackMutex.Unlock()
	pprof.StopCPUProfile()
}

// NotifyOnChanges notifies when events occur (after debouncing)
func (mon *Monitor) NotifyOnChanges() {
	debounceTimer := time.NewTimer(notifyInterval)
	changeset := NewEventChangeset()

	var e notify.EventInfo
	var start time.Time
	for {
		select {
		case e = <-mon.channel:
			// receive an event
			event := e.Event()
			path := e.Path()
			if mon.filter == nil || mon.filter(event, path) {
				if changeset.empty() {
					start = time.Now()
					f, _ := os.Create("cpu.prof")
					pprof.StartCPUProfile(f)
					fmt.Print("Change detected...")
				} else {
					fmt.Print(".")
				}
				changeset.Add(event, path)
				debounceTimer.Reset(mon.debounceDuration)
			}

		case <-debounceTimer.C:
			// debounce and fire callback
			if changeset.nonEmpty() {
				go mon.triggerCallbacks(changeset, start, false)
				changeset = NewEventChangeset()
			} else {
				fmt.Println("")
				fmt.Println("...no changes")
			}
		}
	}
}

// Stop cancels the recursive watcher
func (mon *Monitor) Stop() {
	notify.Stop(mon.channel)
}
