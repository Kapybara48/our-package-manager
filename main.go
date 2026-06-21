package main

import (
	"flag"

	confighelper "our-package-manager/config-helper"
	packagehelper "our-package-manager/package-helper"
)

func main() {
	installFlag := flag.String("install", "", "--install git url of package you want to install")

	flag.Parse()

	if *installFlag != "" {
		//install(*installFlag)
		packageConfig := confighelper.PackageConfig{
			URL:            "https://github.com/Kapybara48/our-package-manager",
			GitCloneDepth:  1,
			GitCloneBranch: "main",
			Makefile:       "Makefile",
			MakefileTarget: "Install",
			InstallScript:  "",
		}

		packagehelper.Install(&packageConfig)
	}
}
