// Package makehelper is wrapper around GNU make command
package makehelper

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"our-package-manager/execute"
)

// Receives directory of the project and executes target like "build" or "install".
func MakeTarget(directory string, makefile string, target string) error {
	err := execMake("-C", directory, "-f", makefile, target)
	if err != nil {
		return err
	}
	return nil
}

func execMake(args ...string) error {
	exitCode, err := execute.ExecuteWithOutput("./", "make", args...)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("make returned exit code %d", exitCode)
	}
	return nil
}

func GetMakeTargets(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var targets []string

	regexTarget, err := regexp.Compile("^[a-zA-Z0-9._%+-]+:")
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return nil, err
		}
		line := scanner.Text()

		if regexTarget.MatchString(line) {
			targets = append(targets, strings.Split(line, ":")[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
