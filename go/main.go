package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	// "path/filepath"

	"our/configuration"
	"our/execute"
	"our/git-helper"

	"github.com/fatih/color"
	// "our/make-helper"
)

var (
	// overridden by -ldflags
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// cli flags
var (
	argHelpPtr   = flag.Bool("help", false, "Help, shows this help message")
	argForcePrt  = flag.Bool("force", false, "Force")
	argDebugPtr  = flag.Bool("debug", false, "Turns on debug mode")
	argConfigPtr = flag.String("config", "/etc/our/config.yml", "Path of configuration file")
	argLimitPtr  = flag.String("limit", ".*", "Regex to limit packages by name")

	argRepoURLPtr     = flag.String("git", "", "URL of the Git repository")
	argPackageNamePtr = flag.String("name", "", "Package name")
	argSubfolderPtr   = flag.String("subfolder", "", "Subfolder in the repository")
)

// cli subcommands options
var (
	subcommandOptions map[string]subcommandOption = map[string]subcommandOption{
		"version": {
			helpMessage: "Shows our version",
			function:    cliVersion,
			loadConfig:  false,
		},
		"add": {
			helpMessage: "Adds package to configuration file (be aware we have to parse the file and write it back, the comments in the file wont survive this action)",
			function:    cliAdd,
			loadConfig:  true,
		},
		"update": {
			helpMessage: "Updates all packages and installs new from configuration file",
			function:    cliUpdate,
			loadConfig:  true,
		},
	}
)

type subcommandOption struct {
	helpMessage string
	loadConfig  bool
	function    func() error
}

func main() {

	if len(os.Args) < 2 {
		cliUsage()
		os.Exit(1)
	}

	flag.CommandLine.Parse(os.Args[2:])
	subcommand, ok := subcommandOptions[os.Args[1]]
	if !ok {
		cliUsage()
		os.Exit(1)
	}

	if subcommand.loadConfig {
		var err error
		err = configuration.LoadConfig(*argConfigPtr)
		if err != nil {
			switch {
			case errors.Is(err, os.ErrNotExist):
				log.Fatalf("Config file \"%s\" does not exist.", *argConfigPtr)
			case errors.Is(err, os.ErrPermission):
				log.Fatalf("Config file \"%s\" incorect permissions.", *argConfigPtr)
			default:
				log.Fatalln(err.Error())
			}

		}
	}

	err := subcommand.function()
	if err != nil {
		log.Fatalln(err)
	}

	// err = makehelper.MakeTarget(filepath.Join(gitRepo.Directory, *subfolder), "Makefile", )
	// if err != nil {
	// 	log.Fatalf("Failed to execute make target: %s", err)
	// }

}

func cliUpdate() error {
	err := configuration.LoadOurState()
	if err != nil {
		log.Printf("error loading packages states: %s\n", err)
		return nil
	}

	limitRegex := regexp.MustCompile(*argLimitPtr)
packageLoop:
	for _, v := range configuration.OurConfigLoaded.PackageConfigs {
		if limitRegex.MatchString(v.Name) == false {
			continue packageLoop
		}
		fmt.Printf(color.YellowString("Checking package %s\n"), v.Name)

		gitRepo := githelper.NewGitRepository(v.GitUrl)

		fmt.Println(color.HiMagentaString("Cloning repository"))

		err := gitRepo.Clone(1, "")
		defer gitRepo.DeleteLocalClone()
		if err != nil {
			fmt.Printf(color.RedString("Failed to clone repository: %s\n"), err)
			continue
		}

		tags, err := gitRepo.GetRemoteTags()
		if err != nil {
			fmt.Printf(color.RedString("Failed to list tags from remote repositroy: %s\n"), err)
			continue
		}

		latestVersionTag := v.GetCompatibleTag(tags)
		if latestVersionTag == "" {
			fmt.Println(color.RedString("Repository does not have tag matrching sleected regex"))
		}

		installedVersion := configuration.GetPackageState(v.Name).InstalledVersion
		if latestVersionTag == installedVersion {
			fmt.Printf(color.GreenString("Latest version of %s, is already installed (%s)\n"), v.Name, latestVersionTag)
			continue
		}

		if installedVersion != "" {
			fmt.Printf(color.HiBlueString("Installed version: %s\n"), installedVersion)
			fmt.Printf(color.GreenString("Updating to: %s\n"), latestVersionTag)
		}

		fmt.Println(color.HiMagentaString("Fetching latest tag matching selected regex"))
		err = gitRepo.FetchTag(latestVersionTag)
		if err != nil {
			fmt.Printf(color.RedString("problem with fetching latest release tag: %s\n"), err)
			continue
		}

		fmt.Print(color.HiMagentaString("Switching to latest tag\n"))
		err = gitRepo.SwitchBranch(latestVersionTag)
		if err != nil {
			fmt.Printf(color.RedString("problem with switching to latest release tag: %s\n"), err)
			continue
		}

		fmt.Println(color.HiMagentaString("Executing installation script"))
		_, exitCode, err := execute.Execute(exec.Cmd{Dir: filepath.Join(gitRepo.Directory, v.PackageSubfolder), Path: v.GetFullPathInstallScript(gitRepo.Directory)}, false, *argDebugPtr)
		if err != nil || exitCode != 0 {
			fmt.Printf(color.RedString("error runing install script exitcode:%d  error: %s\n"), exitCode, err)
			continue
		}

		configuration.SetPackageState(v.Name, configuration.PackageState{InstalledVersion: latestVersionTag, Name: v.Name})
	}
	return nil
}

func cliAdd() error {
	if configuration.OurConfigLoaded.GetPackageByName(*argPackageNamePtr) != nil {
		log.Fatalln("name of the package has to be unique")
	}

	configuration.OurConfigLoaded.PackageConfigs = append(configuration.OurConfigLoaded.PackageConfigs, configuration.PackageConfig{Name: *argPackageNamePtr, PackageSubfolder: *argSubfolderPtr, GitUrl: *argRepoURLPtr})
	err := configuration.OurConfigLoaded.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func cliUsage() error {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [subcommand] [flags]\n\n", color.HiBlueString(os.Args[0]))

	fmt.Fprintln(flag.CommandLine.Output(), "Subcommands:")
	for i, subcommand := range subcommandOptions {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s\n\t%s\n", color.HiGreenString(i), color.RedString(subcommand.helpMessage))
	}
	fmt.Fprintln(flag.CommandLine.Output(), "\nFlags:")
	flag.PrintDefaults()
	return nil
}

func cliVersion() error {
	fmt.Printf("%s %s (Commit: %s, Build Date: %s)\n", color.HiBlueString("our"), color.HiGreenString(Version), Commit, BuildTime)
	return nil
}
