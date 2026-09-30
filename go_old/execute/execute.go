package execute

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecuteGetOutput(workingDir string, command string, args ...string) (string, int, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir
	err := cmd.Start()
	if err != nil {
		return "", 0, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", 0, err
	}

	commandOutput := ""

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()

		err := scanner.Err()
		if err != nil {
			return "", 0, err
		}

		commandOutput = commandOutput + m
	}

	return commandOutput, cmd.ProcessState.ExitCode(), nil
}

func ExecuteWithOutput(workingDir string, command string, args ...string) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, err
	}

	stderrErr := make(chan error)
	stdoutErr := make(chan error)

	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	go outputToStd(stderr, stderrErr)
	go outputToStd(stdout, stdoutErr)

	err = cmd.Wait()
	if err != nil {
		return 0, nil
	}

	err = <-stderrErr
	if err != nil {
		return 0, fmt.Errorf("error while printing stdErr %s", err)
	}
	err = <-stdoutErr
	if err != nil {
		return 0, fmt.Errorf("error while printing stdOut %s", err)
	}

	return cmd.ProcessState.ExitCode(), nil
}

func outputToStd(pipe io.ReadCloser, errorChannel chan error) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		err := scanner.Err()
		if err != nil {
			errorChannel <- err
			return
		}
	}
	errorChannel <- nil
}
