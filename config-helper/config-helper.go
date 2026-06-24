package confighelper

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	DefaultPackageConfig PackageConfig
}

type PackageConfig struct {
	Name           string
	URL            string
	GitCloneDepth  int
	GitCloneBranch string
	Makefile       string
	MakefileTarget string
	InstallScript  string
}

func (p *PackageConfig) SaveConfig() error {
	data, err := toml.Marshal(p)
	if err != nil {
		return err
	}
	err = os.WriteFile("file.toml", data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func ReadPackageConfig(packageConfigPath string) (*PackageConfig, error) {
	var packageConfig *PackageConfig
	data, err := os.ReadFile(packageConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error reading package config file %s", err)
	}
	err = toml.Unmarshal(data, packageConfig)
	if err != nil {
		return nil, fmt.Errorf("error decoding package config file %s", err)
	}
	return packageConfig, nil
}
