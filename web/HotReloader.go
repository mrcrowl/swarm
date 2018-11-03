package web

import (
	"swarm/bundle"
	"swarm/monitor"
	"swarm/source"
)

// HotReloader is responsible for managing hot reloads
type HotReloader struct {
	server    *Server
	workspace *source.Workspace
	moduleSet *bundle.ModuleSet
}

// NewHotReloader creates a new hot reload manager
func NewHotReloader(server *Server, workspace *source.Workspace, moduleSet *bundle.ModuleSet) *HotReloader {
	return &HotReloader{
		server,
		workspace,
		moduleSet,
	}
}

// NotifyReload sends a message to the client page to reload
func (hot *HotReloader) NotifyReload(changes *monitor.EventChangeset) {
	if !hot.server.IsHotReloadEnabled() {
		return
	}

	if changes != nil {
		if changes.SkipHotReload() {
			return
		}

		if changes.HasSingleExt(".css") {
			// css-only reload
			seenFiles := make(map[string]bool)
			for _, change := range changes.Changes() {
				// dedupe: only reload each file once
				if _, seen := seenFiles[change.AbsoluteFilepath()]; seen {
					continue
				}

				seenFiles[change.AbsoluteFilepath()] = true
				if relativePath, ok := hot.workspace.ToRelativePath(change.AbsoluteFilepath()); ok {
					if file := hot.moduleSet.FindFileByPath(relativePath); file != nil {
						cssContent := file.RawContents().(*source.CSSFileContents).RawCSSContent()
						hot.server.TriggerCSSReload(relativePath, cssContent)
					}
				}
			}

			return
		}
	}

	hot.server.TriggerFullReload()
}
