package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	if err := validateCommandAndArgs(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Security error: %v\n", err)
		return 1
	}

	// Безопасный запуск команды
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

func validateCommandAndArgs(cmd []string) error {
	if len(cmd) == 0 {
		return errors.New("empty command")
	}

	if err := validateCommand(cmd[0]); err != nil {
		return err
	}

	// Проверяем аргументы на опасные символы.
	for i, arg := range cmd {
		if i == 0 {
			continue
		}
		if err := validateArgument(arg); err != nil {
			return fmt.Errorf("invalid argument %d: %w", i, err)
		}
	}

	return nil
}

func isAllowedCommand(cmd string) bool {
	allowedCommands := map[string]bool{
		"echo": true,
		"sh":   true,
		"bash": true,
		"go":   true,
		"ls":   true,
	}

	baseCmd := filepath.Base(cmd)
	return allowedCommands[baseCmd]
}

func validateCommand(cmd string) error {
	// Проверяем опасные символы в команде.
	if strings.ContainsAny(cmd, "&|;`$()<>{}[]*?~!\\") {
		return fmt.Errorf("command contains dangerous characters: %s", cmd)
	}

	// Проверяем абсолютные пути к системным директориям.
	if strings.HasPrefix(cmd, "/bin/") ||
		strings.HasPrefix(cmd, "/usr/bin/") ||
		strings.HasPrefix(cmd, "/sbin/") ||
		strings.HasPrefix(cmd, "/usr/sbin/") ||
		strings.Contains(cmd, ":\\") ||
		strings.Contains(cmd, ":/") {
		return fmt.Errorf("absolute paths are not allowed: %s", cmd)
	}

	// Проверяем белый список команд.
	if !isAllowedCommand(cmd) {
		return fmt.Errorf("command not allowed: %s", cmd)
	}

	return nil
}

func validateArgument(arg string) error {
	// Запрещаем аргументы с потенциально опасными символами.
	if strings.ContainsAny(arg, "&|;`$()<>{}[]*?~!\\") {
		return fmt.Errorf("argument contains dangerous characters: %s", arg)
	}

	// Запрещаем аргументы, которые могут быть инъекциями.
	if strings.Contains(arg, "&&") ||
		strings.Contains(arg, "||") ||
		strings.Contains(arg, ";") ||
		strings.Contains(arg, "`") ||
		strings.Contains(arg, "$(") {
		return fmt.Errorf("argument may contain command injection: %s", arg)
	}

	return nil
}
