package packagehelper

import (
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
	gitRepo := githelper.NewGitRepositoryClone(url, 1)
	err := gitRepo.Clone()
	if err != nil {
		return nil, err
	}
	defer gitRepo.DeleteRepository()

	remoteRepoConfigPath := filepath.Join(gitRepo.Directory, "our-info.toml")
	localPackageConfigPath := filepath.Join("/etc/our/packages/", githelper.GetRepositoryNameFromURL(url))

	if fileExists(localPackageConfigPath) {
		packageConfig, err := confighelper.GetPackageConfig(localPackageConfigPath)
		if err != nil {
			return nil, err
		}
		return packageConfig, nil
	}

	if fileExists(remoteRepoConfigPath) {
		packageConfig, err := confighelper.GetPackageConfig(remoteRepoConfigPath)
		if err != nil {
			return nil, err
		}
		err = packageConfig.SaveConfig()
		if err != nil {
			return nil, err
		}
		return packageConfig, nil
	}

	//default
	packageConfig := confighelper.PackageConfig{
		Name:           gitRepo.Name,
		URL:            url,
		GitCloneDepth:  1,
		Makefile:       "Makefile",
		MakefileTarget: "install",
	}

	return &packageConfig, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}
