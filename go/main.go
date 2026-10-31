package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	// "path/filepath"

	"our/configuration"
	"our/execute"
	"our/git-helper"

	"github.com/fatih/color"
	// "our/make-helper"
)

func main() {
	apply := flag.Bool("apply", false, "Applys config file")
	add := flag.Bool("add", false, "Add package")

	configPath := flag.String("config", "", "Path of configuration file")

	debug := flag.Bool("debug", false, "Turns on debug mode")
	repoURL := flag.String("git", "", "URL of the Git repository")
	packageName := flag.String("name", "", "Package name")
	subfolder := flag.String("subfolder", "", "Subfolder in the repository")

	flag.Parse()

	if *configPath != "" {
		err := configuration.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Error loading config file %s", err)
		}
	}

	if *add == true {
		if configuration.OurConfigLoaded.GetPackageByName(*packageName) != nil {
			log.Fatalln("name of the package has to be unique")
		}

		configuration.OurConfigLoaded.PackageConfigs = append(configuration.OurConfigLoaded.PackageConfigs, configuration.PackageConfig{Name: *packageName, PackageSubfolder: *subfolder, GitUrl: *repoURL})
		err := configuration.OurConfigLoaded.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if *apply == true {
		err := configuration.LoadOurState()
		if err != nil {
			log.Printf("error loading packages states: %s\n", err)
			return
		}
		for _, v := range configuration.OurConfigLoaded.PackageConfigs {
			fmt.Printf(color.YellowString("Checking package %s\n"), v.Name)

			gitRepo := githelper.NewGitRepository(v.GitUrl)

			fmt.Println("Cloning repository")
			err := gitRepo.Clone(0, "")
			defer gitRepo.DeleteLocalClone()
			if err != nil {
				fmt.Printf(color.RedString("Failed to clone repository: %s\n"), err)
				continue
			}
			// err = gitRepo.FetchTags()
			// if err != nil {
			// 	fmt.Printf("Failed to fetch tags: %s\n", err)
			// 	continue
			// }

			tags, err := gitRepo.GetTags()
			if err != nil {
				fmt.Printf(color.RedString("Unable to get tags: %s\n"), err)
				continue
			}

			latestVersionTag := v.GetCompatibleTag(tags)
			if latestVersionTag == "" {
				fmt.Print(color.RedString("Git repository does not have tag matching selected regex\n"))
				continue
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

			fmt.Printf("Switching branch to latest tag\n")
			err = gitRepo.SwitchBranch(latestVersionTag)
			if err != nil {
				fmt.Print(color.RedString("problem with switching to latest release tag\n"))
				continue
			}

			fmt.Println("Executing installation script")
			_, exitCode, err := execute.Execute(exec.Cmd{Dir: filepath.Join(gitRepo.Directory, v.PackageSubfolder), Path: v.GetFullPathInstallScript(gitRepo.Directory)}, false, *debug)
			if err != nil || exitCode != 0 {
				fmt.Printf(color.RedString("error runing install script exitcode:%d  error: %s\n"), exitCode, err)
				continue
			}

			configuration.SetPackageState(v.Name, configuration.PackageState{InstalledVersion: latestVersionTag, Name: v.Name})
		}
	}

	// err = makehelper.MakeTarget(filepath.Join(gitRepo.Directory, *subfolder), "Makefile", )
	// if err != nil {
	// 	log.Fatalf("Failed to execute make target: %s", err)
	// }

	fmt.Println("Installation completed successfully.")
}
