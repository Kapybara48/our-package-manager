package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	confighelper "our-package-manager/config-helper"
	packagehelper "our-package-manager/package-helper"
)

func main() {
	installFlag := flag.String("install", "", "git url of package you want to install")
	updateFlag := flag.Bool("update", false, "updates all installed packages")
	listFlag := flag.Bool("list", false, "lists installed packages")
	removeFlag := flag.String("remove", "", "uninstalls package and removes package config")

	flag.Parse()

	if *updateFlag || *listFlag {
		packageConfigFiles, err := os.ReadDir("/etc/our/packages/")
		if err != nil {
			panic(err)
		}

		for _, p := range packageConfigFiles {

			packageConfig, err := confighelper.ReadPackageConfig(filepath.Join("/etc/our/packages/", p.Name()))
			if err != nil {
				panic(err)
			}

			if *updateFlag {
				err = packagehelper.Install(packageConfig)
				if err != nil {
					panic(err)
				}
			}

			if *listFlag {
				fmt.Printf("%s: \"%s\"\n", packageConfig.Name, packageConfig.URL)
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

	if *removeFlag != "" {
		packageConfig, err := packagehelper.GetPackageConfig(*removeFlag)
		if err != nil {
			panic(err)
		}
		err = packagehelper.Uninstall(packageConfig)
		if err != nil {
			panic(err)
		}

	}
}
