package sandbox

import "github.com/spf13/cobra"

// FileReadPermissions allows reading files at the returned paths.
type FileReadPermissions func() []string

type PermissionRegistry interface {
	RegisterFileReadPermission(cmd *cobra.Command, p FileReadPermissions)
}

type Sandbox interface {
	PermissionRegistry
	Apply(cmd *cobra.Command) error
}
