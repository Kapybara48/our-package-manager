package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"our/git-helper"
	"our/make-helper"
)

func main() {
	repoURL := flag.String("repo", "", "URL of the Git repository")
	target := flag.String("target", "install", "Makefile target to execute")
	subfolder := flag.String("subfolder", "", "Subfolder in the repository")
	flag.Parse()

	if *repoURL == "" {
		log.Fatal("Repository URL must be provided")
	}

	gitRepo := githelper.NewGitRepository(*repoURL)

	err := gitRepo.Clone(1, "")
	if err != nil {
		log.Fatalf("Failed to clone repository: %s", err)
	}

	err = makehelper.MakeTarget(filepath.Join(gitRepo.Directory, *subfolder), "Makefile", *target)
	if err != nil {
		log.Fatalf("Failed to execute make target: %s", err)
	}

	fmt.Println("Installation completed successfully.")
}
