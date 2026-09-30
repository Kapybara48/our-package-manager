package execute

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Execute(cmd exec.Cmd, getOutput bool, writeOutputToStd bool) (string, int, error) {
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

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return "", 0, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", 0, err
	}

	err = cmd.Start()
	if err != nil {
		return "", 0, err
	}

	stdoutErr := make(chan error)
	stderrErr := make(chan error)
	stdoutString := make(chan string)
	stderrString := make(chan string)

	go outputToStdAndString(stderr, getOutput, writeOutputToStd, stderrString, stderrErr)
	go outputToStdAndString(stdout, getOutput, writeOutputToStd, stdoutString, stdoutErr)

	err = cmd.Wait()
	if err != nil {
		return "", 0, nil
	}

	err = <-stderrErr
	if err != nil {
		return "", 0, fmt.Errorf("error while printing stdErr %s", err)
	}
	err = <-stdoutErr
	if err != nil {
		return "", 0, fmt.Errorf("error while printing stdOut %s", err)
	}

	return commandOutput, cmd.ProcessState.ExitCode(), nil
}

// outputToStdAndString is used to write output to std or read it into string and return it, or if it is usefull, both
func outputToStdAndString(pipe io.ReadCloser, getOutput bool, writeOutputToStd bool, output chan string, errorChannel chan error) {
	strBuilder := strings.Builder{}
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		m := scanner.Text()
		err := scanner.Err()
		if err != nil {
			errorChannel <- err
			output <- strBuilder.String()
			return
		}
		if writeOutputToStd {
			fmt.Println(m)
		}
		if getOutput {
			strBuilder.WriteString(m + "\n")
		}
	}
	output <- strBuilder.String()
	errorChannel <- nil
}
