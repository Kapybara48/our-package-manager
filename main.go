package main

import (
	"flag"
	"os"
	"path/filepath"

	confighelper "our-package-manager/config-helper"
	packagehelper "our-package-manager/package-helper"
)

func main() {
	installFlag := flag.String("install", "", "--install git url of package you want to install")
	updateFlag := flag.Bool("update", false, "Updates all installed packages")

	flag.Parse()

	if *updateFlag {
		packageConfigFiles, err := os.ReadDir("/etc/our/packages/")
		if err != nil {
			panic(err)
		}

		for _, p := range packageConfigFiles {

			packageConfig, err := confighelper.ReadPackageConfig(filepath.Join("/etc/our/packages/", p.Name()))
			if err != nil {
				panic(err)
			}
			err = packagehelper.Install(packageConfig)
			if err != nil {
				panic(err)
			}
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
