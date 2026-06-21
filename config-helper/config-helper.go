package confighelper

type PackageConfig struct {
	URL            string
	GitCloneDepth  int
	GitCloneBranch string
	Makefile       string
	MakefileTarget string
	InstallScript  string
}
