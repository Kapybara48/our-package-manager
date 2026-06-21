package confighelper

type PackageConfig struct {
	Url            string
	GitCloneDepth  string
	GitCloneBranch string
	Makefile       string
	MakefileTarget string
	InstallScript  string
}
