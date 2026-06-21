// Package githelper
package githelper

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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

	cmd := exec.Command("git", "clone", generateFolderName(g.URL), "/tmp/our-package-manager")
	cmd.Args = append(cmd.Args, args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	return nil
}

func generateFolderName(url string) string {
	parts := strings.Split(url, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return "/tpm/" + repoName
}

func generateRandomString(length int) string {
}
