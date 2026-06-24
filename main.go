package main

import (
	"flag"

	packagehelper "our-package-manager/package-helper"
)

func main() {
	installFlag := flag.String("install", "", "--install git url of package you want to install")

	flag.Parse()

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
	}
}
