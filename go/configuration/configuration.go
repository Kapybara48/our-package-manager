package configuration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/muesli/go-app-paths"
)

type PackageConfig struct {
	Name             string   `toml:"Name"`
	GitUrl           string   `toml:"GitUrl"`
	PackageSubfolder string   `toml:"PackageSubfolder"`
	Branch           string   `toml:"Branch"`
	InstallScript    string   `toml:"InstallScript"`
	ReleaseTagRegex  string   `toml:"ReleaseTagRegex"`
	PackageFolders   []string `toml:"PackageFolders"`
}

type PackageState struct {
	Name             string `toml:"Name"`
	InstalledVersion string `toml:"InstalledVersion"`
}

type AppState struct {
	PackagesStates []PackageState
}

type AppConfig struct {
	PackageConfigs       []PackageConfig `yaml:"Package" json:"Package" toml:"Package"`
	DefualtInstallScript string          `yaml:"DefaultInstallScript" json:"DefaultInstallScript" toml:"DefaultInstallScript"`
}

var ourConfigLoadedPath string
var OurConfigLoaded *AppConfig = &AppConfig{}

var ourStatePath string
var OurState *AppState = &AppState{}

func LoadConfig(path string) error {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(configBytes, OurConfigLoaded)
	if err != nil {
		return err
	}
	ourConfigLoadedPath = path

	return nil
}

func (c *AppConfig) WriteConfig() error {
	configBytes, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(ourConfigLoadedPath, configBytes, 0644)
	return nil
}

func SetPackageState(name string, ps PackageState) {
	exists := false
	for i := range OurState.PackagesStates {
		if OurState.PackagesStates[i].Name == name {
			OurState.PackagesStates[i] = ps
			exists = true
		}
	}

	if exists == false {
		OurState.PackagesStates = append(OurState.PackagesStates, ps)
	}

	err := WriteOurState()
	if err != nil {
		panic("cannot write our state into data file")
	}
}

func GetPackageState(name string) PackageState {
	for i := range OurState.PackagesStates {
		if OurState.PackagesStates[i].Name == name {
			return OurState.PackagesStates[i]
		}
	}
	return PackageState{}
}

func LoadOurState() error {
	scope := gap.NewScope(gap.System, "our")
	statePath, err := scope.DataPath("OurState.toml")
	if err != nil {
		return err
	}

	packageStateBytes, err := os.ReadFile(statePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !os.IsNotExist(err) {
		err = toml.Unmarshal(packageStateBytes, &OurState)
		if err != nil {
			return err
		}
	}

	ourStatePath = statePath

	return nil
}

func WriteOurState() error {
	packageStatesBytes, err := toml.Marshal(OurState)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(ourStatePath), 0o755); err != nil {
		return err
	}

	err = os.WriteFile(ourStatePath, packageStatesBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *AppConfig) GetPackageByName(name string) *PackageConfig {
	for i := range c.PackageConfigs {
		if c.PackageConfigs[i].Name == name {
			return &c.PackageConfigs[i]
		}
	}
	return nil
}

func (c *PackageConfig) GetCompatibleTag(tags []string) string {
	regexCompiled, err := regexp.Compile(c.ReleaseTagRegex)
	if err != nil {
		return ""
	}

	for _, v := range tags {
		if regexCompiled.Match([]byte(v)) {
			return v
		}
	}

	return ""
}

func (c *PackageConfig) GetFullPathInstallScript(gitRepoPath string) string {
	installScriptFullPath := c.InstallScript
	if !strings.HasPrefix(installScriptFullPath, "/") {
		if after, found := strings.CutPrefix(installScriptFullPath, "git:"); found == true {
			installScriptFullPath = filepath.Join(gitRepoPath, after)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				return ""
			}
			installScriptFullPath = filepath.Join(wd, installScriptFullPath)
		}
	}

	return installScriptFullPath
}
