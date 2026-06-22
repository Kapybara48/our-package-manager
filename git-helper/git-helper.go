// Package githelper
package githelper

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	confighelper "our-package-manager/config-helper"
	"our-package-manager/execute"
)

type GitRepository struct {
	URL       string
	Depth     int
	Branch    string
	Directory string
	Name      string
}

func NewGitRepository(packageConfig confighelper.PackageConfig) *GitRepository {
	git := GitRepository{}

	git.URL = packageConfig.URL
	git.Depth = packageConfig.GitCloneDepth
	git.Branch = packageConfig.GitCloneBranch
	git.Name = GetRepositoryNameFromURL(git.URL)
	git.Directory = generateFolderName(git.URL)

	return &git
}

func (g *GitRepository) Clone() error {
	var args []string

	args = append(args, "clone", g.URL, g.Directory)

	if g.Depth != 0 {
		args = append(args, "--depth", strconv.Itoa(g.Depth))
	}

	if g.Branch != "" {
		args = append(args, "--branch", g.Branch)
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

func (g *GitRepository) SwitchBranch() error {
	exitCode, err := execute.ExecuteWithOutput(g.Directory, "git", "switch", g.Branch)
	if err != nil {
		return fmt.Errorf("error switching branch %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("switching branch returned exit code %d", exitCode)
	}
	return nil
}

func (g *GitRepository) DeleteRepository() error {
	err := os.RemoveAll(g.Directory)
	if err != nil {
		return fmt.Errorf("Error removing directory: %v", err)
	}
	return nil
}

func generateFolderName(repositoryName string) string {
	return "/tmp/" + repositoryName + "-" + generateRandomString(10)
}

func GetRepositoryNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return repoName
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
