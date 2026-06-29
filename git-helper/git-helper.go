// Package githelper
package githelper

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"our-package-manager/execute"
)

type GitRepository struct {
	URL       string
	Directory string // if Directory is empty, it means, the repo was not cloned yet
}

func NewGitRepository(url string) *GitRepository {
	git := GitRepository{}

	git.URL = url
	name := git.GetName()
	git.Directory = generateFolderName(name)

	return &git
}

func (g *GitRepository) Clone(depth int, branch string) error {
	var args []string

	args = append(args, "clone", g.URL, g.Directory)

	if depth != 0 {
		args = append(args, "--depth", strconv.Itoa(depth))
	}

	if branch != "" {
		args = append(args, "--branch", branch)
	}

	exitCode, err := execute.ExecuteWithOutput(".", "git", args...)
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("cloning repository returned exit code %d", exitCode)
	}

	return nil
}

func (g *GitRepository) SwitchBranch(branch string) error {
	exitCode, err := execute.ExecuteWithOutput(g.Directory, "git", "switch", branch)
	if err != nil {
		return fmt.Errorf("error switching branch %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("switching branch returned exit code %d", exitCode)
	}
	return nil
}

func (g *GitRepository) DeleteLocalClone() error {
	err := os.RemoveAll(g.Directory)
	if err != nil {
		return fmt.Errorf("error removing git directory: %s", err)
	}

	g.Directory = ""

	return nil
}

func (g *GitRepository) GetName() string {
	if g.URL == "" || g.URL == "." {
		cwd, err := os.Getwd()
		if err == nil {
			g.URL = cwd
			return g.GetName()
		}
		return "local-repo"
	}

	parts := strings.Split(g.URL, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return repoName
}

func generateFolderName(repositoryName string) string {
	return filepath.Join(os.TempDir(), repositoryName+"-"+generateRandomString(10))
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
