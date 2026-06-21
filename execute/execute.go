package execute

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecuteWithOutput(command string, args ...string) (error, int) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, 0
	}

	stderrErr := make(chan error)
	stdoutErr := make(chan error)
	err = cmd.Start()
	if err != nil {
		return err, 0
	}

	go outputToStd(stderr, stderrErr)
	go outputToStd(stdout, stdoutErr)

	cmd.Wait()

	err = <-stderrErr
	if err != nil {
		return err, 0
	}
	err = <-stdoutErr
	if err != nil {
		return err, 0
	}
	return nil, cmd.ProcessState.ExitCode()
}

func outputToStd(pipe io.ReadCloser, errorChannel chan error) {
	defer pipe.Close()

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
