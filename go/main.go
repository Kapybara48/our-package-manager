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

	packageLoop:
		for _, v := range configuration.OurConfigLoaded.PackageConfigs {
			fmt.Printf(color.YellowString("Checking package %s\n"), v.Name)

			gitRepo := githelper.NewGitRepository(v.GitUrl)

			fmt.Println(color.HiMagentaString("Cloning repository"))
			latestVersionTag := ""

			err := gitRepo.Clone(1, "")
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

			lastRepoDepth := 1

			for latestVersionTag == "" {
				fmt.Println(color.HiMagentaString("Deepening repository by 30 commits"))

				err := gitRepo.Deepen(30)
				if err != nil {
					fmt.Printf(color.RedString("Failed to deepen repository: %s\n"), err)
					continue packageLoop
				}

				repoDepth, err := gitRepo.CountCommits()
				if err != nil {
					fmt.Printf(color.RedString("Faled to get commit count of repo: %s"), err)
				}

				if repoDepth == lastRepoDepth {
					fmt.Println(color.RedString("Git repository does not have tag matching selected regex"))
					continue packageLoop
				}
				lastRepoDepth = repoDepth

				tags, err := gitRepo.GetTags()
				if err != nil {
					fmt.Printf(color.RedString("Unable to get tags: %s\n"), err)
					continue packageLoop
				}

				latestVersionTag = v.GetCompatibleTag(tags)
				if latestVersionTag == "" {
					fmt.Println(color.HiMagentaString("Git repository does not have tag matching selected regex (in current depth)"))
				}
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

			fmt.Print(color.HiMagentaString("Switching branch to latest tag\n"))
			err = gitRepo.SwitchBranch(latestVersionTag)
			if err != nil {
				fmt.Print(color.RedString("problem with switching to latest release tag\n"))
				continue
			}

			fmt.Println(color.HiMagentaString("Executing installation script"))
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

}
