package makehelper

import (
	"fmt"

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
	exitCode, err := execute.ExecuteWithOutput("make", args...)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("make returned exit code %d", exitCode)
	}
	return nil
}
