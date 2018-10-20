package dep

import (
	"gospm/dep"
	"log"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

// WorkspaceWatcher is used to recursively watch for file changes within a workspace
type WorkspaceWatcher struct {
	workspace *dep.Workspace
	channel   chan notify.EventInfo
}

// NewWorkspaceWatcher creates a new WorkspaceWatcher
func NewWorkspaceWatcher(workspace *dep.Workspace) *WorkspaceWatcher {
	channel := make(chan notify.EventInfo, 2048)

	rootRecursivePattern := filepath.Join(workspace.RootPath(), "./...")
	if err := notify.Watch(rootRecursivePattern, channel, notify.All); err != nil {
		log.Fatal(err)
	}

	return &WorkspaceWatcher{
		workspace,
		channel,
	}
}

// C retursn the channel which notifies about events
func (wsw *WorkspaceWatcher) C() chan notify.EventInfo {
	return wsw.channel
}

// Stop cancels the recursive watcher
func (wsw *WorkspaceWatcher) Stop() {
	notify.Stop(wsw.channel)
}
