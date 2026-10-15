package sandbox

import (
	"errors"

	"github.com/spf13/cobra"
)

type Landlock struct {
	readFiles map[*cobra.Command]FileReadPermissions
}

func (l *Landlock) RegisterFileReadPermission(cmd *cobra.Command, p FileReadPermissions) {
	if l.readFiles == nil {
		l.readFiles = make(map[*cobra.Command]FileReadPermissions)
	}
	l.readFiles[cmd] = p
}

var UnsupportedPlatformError error = errors.New("this platform does not support landlock")
