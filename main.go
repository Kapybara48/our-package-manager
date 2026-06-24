package main

import (
	"flag"
	"fmt"
	"os"

	githelper "our-package-manager/git-helper"
	packagehelper "our-package-manager/package-helper"
)

func main() {
	gitRepo := githelper.NewGitRepositoryClone("https://github.com/bedrixh/ssl-manager.git", 1)
	fmt.Printf("(%s)", gitRepo)
	os.Exit(0)
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
