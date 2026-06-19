package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

func Clone(url string, destination string, depth int, branch string) error {
	cmd := exec.Command("git", "clone", "--depth", strconv.Itoa(depth), "--branch", branch, url, destination)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning git repository %s", err)
	}
	return nil
}
