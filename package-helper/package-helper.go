package packagehelper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	confighelper "our-package-manager/config-helper"
	githelper "our-package-manager/git-helper"
	makehelper "our-package-manager/make-helper"
)

func Install(packageConfig *confighelper.PackageConfig) error {
	git := githelper.NewGitRepository(*packageConfig)
	err := git.Clone()
	if err != nil {
		return err
	}

	err = makehelper.MakeTarget(git.Directory, packageConfig.Makefile, "install")
	if err != nil {
		return err
	}

	log.Printf("\"%s\" installed successfully\n", git.URL)

	defer git.DeleteRepository()
	return nil
}

func GetPackageConfig(url string) (*confighelper.PackageConfig, error) {
	gitRepo := githelper.NewGitRepositoryClone(url)
	err := gitRepo.Clone()
	if err != nil {
		return nil, err
	}

	fileExists(filepath.Join(gitRepo.Directory, ""))
	confighelper.GetPackageConfig()

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
