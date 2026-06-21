// Package git
package githelper

import (
	"fmt"
	"os/exec"
	"strconv"
)

type GitRepository struct {
	URL         string
	Destination string
	Depth       int
	Branch      string
}

func (g *GitRepository) Clone() error {
	cmd := exec.Command("git", "clone", "--depth", strconv.Itoa(g.Depth), "--branch", g.Branch, g.URL, g.Destination)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	return nil
}
