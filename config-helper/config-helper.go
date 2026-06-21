package confighelper

type PackageConfig struct {
	Url            string
	GitCloneDepth  string
	Makefile       string
	MakefileTarget string
	InstallScript  string
}
