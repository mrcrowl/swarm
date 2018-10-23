package monitor

import (
	"fmt"
	"os"
	"swarm/config"
	"swarm/source"
	"swarm/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	dir := testutil.CreateTempDirWithPrefix("TestMonitor")
	ws := source.NewWorkspace(dir)
	mon := NewMonitor(ws, config.NewMonitorConfig([]string{".js"}, 150))

	notifyCount := 0
	eventCount := 0
	go mon.NotifyOnChanges(func(ec *EventChangeset) {
		notifyCount++
		eventCount += ec.count()
	})

	defer os.RemoveAll(dir)

	for i := 0; i < 10; i++ {
		var filename string
		if i < 9 {
			filename = fmt.Sprintf("abcd%d.js", i)
		} else {
			filename = "abcd.ts"
		}
		testutil.WriteTextFile(dir, filename, "hello world")
	}

	assert.Equal(t, 0, notifyCount)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, notifyCount)
	assert.Equal(t, 9, eventCount) // 10 - 1 filtered
}
