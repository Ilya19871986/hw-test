package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	if err := validateCommand(cmd[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Security error: %v\n", err)
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = env.PrepareEnv()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

func isAllowedCommand(cmd string) bool {
	allowedCommands := map[string]bool{
		"echo": true,
		"sh":   true,
		"bash": true,
		"go":   true,
	}

	baseCmd := filepath.Base(cmd)
	return allowedCommands[baseCmd]
}

func validateCommand(cmd string) error {
	if !isAllowedCommand(cmd) {
		return fmt.Errorf("command not allowed: %s", cmd)
	}
	return nil
}
