package main

import (
	"flag"
	"our-package-manager/git"
)

func main() {
	//	nFlag := flag.Int("n", 1234, "help message for flag n")
	//	flag.Parse()
	//	if *nFlag == 1234 {
	//		fmt.Println("n flag used")
	//	}

	//gitUrlFlag := flag.String("git-url", "", "link for git project")
	flag.Parse()
	err := git.Clone("https://github.com/Foxboron/sbctl.git", "/tmp/sbctl-git", 1, "master")
	if err != nil {
		panic(err)
	}
}
