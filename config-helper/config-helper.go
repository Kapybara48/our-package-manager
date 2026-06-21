package confighelper

type PackageConfig struct {
	URL            string
	GitCloneDepth  string
	GitCloneBranch string
	Makefile       string
	MakefileTarget string
	InstallScript  string
}
