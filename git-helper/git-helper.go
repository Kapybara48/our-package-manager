// Package githelper
package githelper

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"

	confighelper "our-package-manager/config-helper"
)

type GitRepository struct {
	URL       string
	Depth     int
	Branch    string
	Directory string
}

func NewGitRepository(packageConfig confighelper.PackageConfig) *GitRepository {
	git := GitRepository{}

	git.URL = packageConfig.URL
	git.Depth = packageConfig.GitCloneDepth
	git.Branch = packageConfig.GitCloneBranch
	git.Directory = generateFolderName(git.URL)

	return &git
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
	return "/tpm/" + repoName + generateRandomString(10)
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	charsetLength := len(charset)
	randomChars := make([]byte, length)

	for i := range length {
		randomChars[i] = charset[rand.Intn(charsetLength)]
	}

	return string(randomChars)
}
