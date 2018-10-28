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

// MimeTypeFromFilename returns the mimetype for a filename based on its filename (for a few selected extensions)
func MimeTypeFromFilename(filename string) string {
	ext := path.Ext(filename)
	switch ext {
	case ".js":
		return "application/javascript"
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	}
	return "text/plain; charset=utf-8"
}
