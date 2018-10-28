// +build !embed

package assets

import (
	"github.com/omeid/go-resources/live"
)

// FS is the exported file system
var FS = live.Dir("..")
