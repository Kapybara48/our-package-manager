package packagehelper

import (
	"log"

	confighelper "our-package-manager/config-helper"
	githelper "our-package-manager/git-helper"
	makehelper "our-package-manager/make-helper"
)

func Install(packageConfig *confighelper.PackageConfig) error {
	git := githelper.GitRepository{URL: packageConfig.Url, Depth: 1}
	err := git.Clone()
	if err != nil {
		panic(err)
	}

	err = makehelper.MakeTarget(git.Directory, packageConfig.Makefile, "install")
	if err != nil {
		return err
	}

	log.Println("\"%s\" installed successfully", git.URL)

	return nil
}
