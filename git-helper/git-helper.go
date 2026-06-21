// Package githelper
package githelper

import (
	"fmt"
	"os/exec"
	"strconv"
)

type GitRepository struct {
	URL    string
	Depth  int
	Branch string
}

func (g *GitRepository) Clone() error {
	var args []string

	if g.Depth != 0 {
		args = append(args, "--depth", strconv.Itoa(g.Depth))
	}

	if g.Branch != "" {
		args = append(args, "--branch", g.Branch)
	}

	cmd := exec.Command("git", "clone", g.URL, "/tmp/our-package-manager")
	cmd.Args = append(cmd.Args, args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	return nil
}
