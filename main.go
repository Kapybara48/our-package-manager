package main

import (
	"flag"

	packagehelper "our-package-manager/package-helper"
)

func main() {
	installFlag := flag.String("install", "", "--install git url of package you want to install")
	selfUpdateFlag := flag.Bool("self-update", false, "Self updates our-package-manager")

	flag.Parse()

	if *selfUpdateFlag {
		packageConfig, err := packagehelper.GetPackageConfig("https://github.com/Kapybara48/our-package-manager.git")
		if err != nil {
			panic(err)
		}
		err = packagehelper.Install(packageConfig)
		if err != nil {
			panic(err)
		}

		return
	}

	if *installFlag != "" {
		//install(*installFlag)
		//		packageConfig := confighelper.PackageConfig{
		//			URL:            *installFlag,
		//			GitCloneDepth:  1,
		//			GitCloneBranch: "main",
		//			Makefile:       "Makefile",
		//			MakefileTarget: "Install",
		//			InstallScript:  "",
		//		}
		packageConfig, err := packagehelper.GetPackageConfig(*installFlag)
		if err != nil {
			panic(err)
		}
		err = packagehelper.Install(packageConfig)
		if err != nil {
			panic(err)
		}
		return
	}
}
