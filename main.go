package main

import (
	"flag"
	"log"

	githelper "our-package-manager/git-helper"
	makehelper "our-package-manager/make-helper"
)

func main() {
	installFlag := flag.String("install", "", "--install git url of package you want to install")

	flag.Parse()

	if *installFlag != "" {
		install(*installFlag)
	}
}

func install(packageUrl string) error {
	git := githelper.GitRepository{URL: packageUrl, Depth: 1}
	err := git.Clone()
	if err != nil {
		panic(err)
	}

	err = makehelper.MakeTarget(git.Directory, "install")
	if err != nil {
		return err
	}

	log.Println("\"%s\" installed successfully", git.URL)

	return nil
}
