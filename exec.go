package main

import (
	"fmt"
	"os"
	"os/exec"
)

func executeCommand(args []string) int {
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: command '%s' not found in PATH\n", args[0])
		return 127
	}

	cmd := exec.Command(cmdPath, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	fmt.Fprintf(os.Stderr, "Error: failed to execute command: %v\n", err)
	return 1
}
