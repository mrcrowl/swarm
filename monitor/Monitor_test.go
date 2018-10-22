package monitor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"swarm/source"
	"testing"
	"time"

	"github.com/rjeczalik/notify"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "TestMonitor")
	if err != nil {
		log.Panicln("Failed to create temp dir")
	}

	ws := source.NewWorkspace(dir)
	filterFn := func(event notify.Event, path string) bool { return !strings.HasSuffix(path, "9.js") }
	mon := NewMonitor(ws, filterFn)

	notifyCount := 0
	eventCount := 0
	go mon.NotifyOnChanges(func(ec *EventChangeset) {
		notifyCount++
		eventCount += ec.count()
	})

	defer os.RemoveAll(dir)

	for i := 0; i < 10; i++ {
		filename := fmt.Sprintf("abcd%d.js", i)
		ioutil.WriteFile(filepath.Join(dir, filename), []byte("hello world"), os.ModePerm)
	}

	assert.Equal(t, 0, notifyCount)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, notifyCount)
	assert.Equal(t, 9, eventCount) // 10 - 1 filtered
}
