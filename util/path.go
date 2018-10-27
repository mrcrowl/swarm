package util

import "path"

// RemoveExtension returns a path without the extension
func RemoveExtension(relativePath string) string {
	ext := path.Ext(relativePath)
	if ext != "" {
		return relativePath[:len(relativePath)-len(ext)]
	}
	return relativePath
}
