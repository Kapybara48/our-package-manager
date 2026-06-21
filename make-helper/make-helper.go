package makehelper

import (
	"fmt"
	"our-package-manager/execute"
)

// Receives directory of the project and target like "build" or "install". 
func MakeTarget(directory string, target string) error {
	err := execMake("-c", directory, target)
	if err != nil {
		return err
	}
	return nil
}

func execMake(args ...string) error {
	err, exitCode := execute.ExecuteWithOutput("make",args...)
	if err != nil {
		return err
	}
	if exitCode != 0{
		return fmt.Errorf("exit code not 0 after running make")
	}
	return nil
}
