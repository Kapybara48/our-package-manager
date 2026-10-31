// Package githelper
package githelper

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"our/execute"
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

	_, exitCode, err := execute.Execute(exec.Cmd{Args: args, Path: "git"}, false, true)
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("cloning repository returned exit code %d", exitCode)
	}

	return nil
}

func (g *GitRepository) SwitchBranch(branch string) error {
	_, exitCode, err := execute.Execute(exec.Cmd{Dir: g.Directory, Path: "git", Args: []string{"switch", branch}}, false, true)
	if err != nil {
		return fmt.Errorf("error switching branch %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("switching branch returned exit code %d", exitCode)
	}
	return nil
}

func (g *GitRepository) FetchTags() error {
	_, exitCode, err := execute.Execute(exec.Cmd{Dir: g.Directory, Path: "git", Args: []string{"fetch", "--tags"}}, false, false)
	if err != nil {
		return fmt.Errorf("error fetching tags %s", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("fetching tags returned exit code %d", exitCode)
	}
	return nil
}

func (g *GitRepository) GetTags() (error, []string) {
	tags := []string{}
	_, exitCode, err := execute.Execute(exec.Cmd{Dir: g.Directory, Path: "git", Args: []string{"tag", "-l", "--sort=-creatordate"}}, true, false)
	if exitCode != 0 {
		return fmt.Errorf("error geting tags: git fetch exited with %d", exitCode), tags
	}

	if err != nil {
		return fmt.Errorf("error geting tags: %w", err), tags
	}

	return nil, tags
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
