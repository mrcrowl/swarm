package monitor

import (
	"github.com/rjeczalik/notify"
)

// EventChangeset is used to collect a set of events
type EventChangeset struct {
	changeIndex map[string]bool
	changes     []*Event
}

// NewEventChangeset creates a new EventChangeset
func NewEventChangeset() *EventChangeset {
	return &EventChangeset{
		changeIndex: make(map[string]bool),
		changes:     nil,
	}
}

// Changes get the list of changes that make up this
func (ec *EventChangeset) Changes() []*Event {
	return ec.changes
}

// Add adds a new event to the set
func (ec *EventChangeset) Add(event notify.Event, path string) bool {
	if isCompositeEvent(event) {
		return false
	}
	key := makeEventKey(event, path)
	if _, exists := ec.changeIndex[key]; !exists {
		ev := NewEvent(path, event)
		ec.changes = append(ec.changes, ev)
		ec.changeIndex[key] = true
		return true
	}

	return false
}

func (ec *EventChangeset) count() int {
	return len(ec.changeIndex)
}

func (ec *EventChangeset) empty() bool {
	return ec.count() == 0
}

func (ec *EventChangeset) nonEmpty() bool {
	return ec.count() > 0
}

func makeEventKey(event notify.Event, path string) string {
	return eventToString(event) + ":" + path
}

func isCompositeEvent(e notify.Event) bool {
	n := int(e)
	return (n-1)&n > 0
}

func eventToString(e notify.Event) string {
	switch e {
	case notify.Create:
		return "C"
	case notify.Write:
		return "W"
	case notify.Remove:
		return "D"
	case notify.Rename:
		return "M"
	default:
		return "?"
	}
}
