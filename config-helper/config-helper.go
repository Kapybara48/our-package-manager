package confighelper

import (
	"os"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
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

func LoadPackageConfig() (*PackageConfig, error) {
	var packageConfig *PackageConfig
	data, err := os.ReadFile("")
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, packageConfig)
	if err != nil {
		return nil, err
	}
	return packageConfig, nil
}
