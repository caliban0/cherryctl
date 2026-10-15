//go:build !linux

package sandbox

import "github.com/spf13/cobra"

func (l *Landlock) Apply(_ *cobra.Command) error {
	return UnsupportedPlatformError
}
