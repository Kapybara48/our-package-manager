package main

import (
	"flag"

	githelper "our-package-manager/git-helper"
)

func main() {
	//	nFlag := flag.Int("n", 1234, "help message for flag n")
	//	flag.Parse()
	//	if *nFlag == 1234 {
	//		fmt.Println("n flag used")
	//	}

	//gitUrlFlag := flag.String("git-url", "", "link for git project")
	flag.Parse()
	git := githelper.GitRepository{URL: "https://github.com/Foxboron/sbctl.git", Depth: 1, Branch: "master"}
	err := git.Clone()
	if err != nil {
		panic(err)
	}
}
