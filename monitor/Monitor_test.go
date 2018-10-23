package monitor

import (
	"fmt"
	"os"
	"strings"
	"swarm/source"
	"swarm/testutil"
	"testing"
	"time"

	"github.com/rjeczalik/notify"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	dir := testutil.CreateTempDirWithPrefix("TestMonitor")
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
		testutil.WriteTextFile(dir, filename, "hello world")
	}

	assert.Equal(t, 0, notifyCount)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, notifyCount)
	assert.Equal(t, 9, eventCount) // 10 - 1 filtered
}
