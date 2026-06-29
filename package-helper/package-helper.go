package packagehelper

import (
	"os"
	"path/filepath"

	confighelper "our-package-manager/config-helper"
	githelper "our-package-manager/git-helper"
	makehelper "our-package-manager/make-helper"
)

func Install(packageConfig *confighelper.PackageConfig) error {
	git := githelper.NewGitRepository(packageConfig.URL)
	err := git.Clone(packageConfig.GitDepth, packageConfig.GitBranch)
	if err != nil {
		return err
	}

	err = makehelper.MakeTarget(git.Directory, packageConfig.Makefile, "install")
	if err != nil {
		return err
	}

	defer git.DeleteLocalClone()
	return nil
}

func Uninstall(packageConfig *confighelper.PackageConfig) error {
	gitRepo := githelper.NewGitRepository(packageConfig.URL)

	err := gitRepo.Clone(packageConfig.GitDepth, packageConfig.GitBranch)
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

	defer gitRepo.DeleteLocalClone()
	return nil

}

func GetPackageConfig(url string) (*confighelper.PackageConfig, error) {
	gitRepo := githelper.NewGitRepository(url)

	localPackageConfigPath := filepath.Join("/etc/our/packages/", gitRepo.GetName()+".toml")
	remoteRepoConfigPath := filepath.Join(gitRepo.Directory, "our-info.toml")

	if fileExists(localPackageConfigPath) {
		packageConfig, err := confighelper.ReadPackageConfig(localPackageConfigPath)
		if err != nil {
			return nil, err
		}
		return packageConfig, nil
	}

	// remote config loading
	err := gitRepo.Clone(1, "")
	if err != nil {
		return nil, err
	}
	defer gitRepo.DeleteLocalClone()

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
		Name:                gitRepo.GetName(),
		URL:                 url,
		GitDepth:            1,
		Makefile:            "Makefile",
		MakeInstallTarget:   "install",
		MakeUninstallTarget: "uninstall",
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
