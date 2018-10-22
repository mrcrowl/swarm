package monitor

import "github.com/rjeczalik/notify"

// Event represents a file-level change
type Event struct {
	path  string // TODO <-- rename to filepath
	event notify.Event
}

// NewEvent creates a new Event instance
func NewEvent(path string, event notify.Event) *Event {
	return &Event{path, event}
}

// AbsoluteFilepath gets the absolute filepath for this Evenet
func (e *Event) AbsoluteFilepath() string {
	return e.path
}
