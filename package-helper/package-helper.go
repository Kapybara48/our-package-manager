package packagehelper

import (
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

	defer git.DeleteRepository()
	return nil
}

func Uninstall(packageConfig *confighelper.PackageConfig) error {
	gitRepo := githelper.NewGitRepository(*packageConfig)

	err := gitRepo.Clone()
	if err != nil {
		return err
	}

	err = makehelper.MakeTarget(gitRepo.Directory, packageConfig.Makefile, packageConfig.MakeUninstallTarget)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join("/etc/our/packages/", packageConfig.Name+".toml"))
	if err != nil {
		return err
	}

	defer gitRepo.DeleteRepository()
	return nil

}

func GetPackageConfig(url string) (*confighelper.PackageConfig, error) {
	gitRepo := githelper.NewGitRepositoryClone(url, 1)
	localPackageConfigPath := filepath.Join("/etc/our/packages/", gitRepo.Name+".toml")
	remoteRepoConfigPath := filepath.Join(gitRepo.Directory, "our-info.toml")

	if fileExists(localPackageConfigPath) {
		packageConfig, err := confighelper.ReadPackageConfig(localPackageConfigPath)
		if err != nil {
			return nil, err
		}
		return packageConfig, nil
	}

	// remote config loading
	err := gitRepo.Clone()
	if err != nil {
		return nil, err
	}
	defer gitRepo.DeleteRepository()

	if fileExists(remoteRepoConfigPath) {
		packageConfig, err := confighelper.ReadPackageConfig(remoteRepoConfigPath)
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
		Name:              gitRepo.Name,
		URL:               url,
		GitCloneDepth:     1,
		Makefile:          "Makefile",
		MakeInstallTarget: "install",
	}
	err = packageConfig.SaveConfig()
	if err != nil {
		return nil, err
	}

	return &packageConfig, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
