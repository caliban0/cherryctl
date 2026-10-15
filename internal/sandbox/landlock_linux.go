//go:build linux

package sandbox

import (
	"fmt"
	"path/filepath"

	"github.com/landlock-lsm/go-landlock/landlock"
	"github.com/spf13/cobra"
)

const certPath = "/etc/ssl/certs/ca-certificates.crt"

func (l *Landlock) Apply(cmd *cobra.Command) error {
	var paths []string
	for c := cmd; c != nil; c = c.Parent() {
		if f := l.readFiles[c]; f != nil {
			fmt.Println(paths)
			paths = append(paths, l.readFiles[c]()...)
		}
	}

	paths = append(paths, "/etc/resolv.conf")

	resolvedCertPath, err := filepath.EvalSymlinks(certPath)
	if err != nil {
		return err
	}

	return landlock.V10.BestEffort().RestrictPaths(
		landlock.ROFiles(paths...),
		landlock.RODirs("/etc/ssl/certs"),
		landlock.ROFiles(resolvedCertPath),
	)
}
