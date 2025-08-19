package main

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("successful command", func(t *testing.T) {
		cmd := []string{"echo", "hello world"}
		env := Environment{}

		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
	})

	t.Run("command with environment", func(t *testing.T) {
		cmd := []string{"sh", "-c", "echo $TEST_VAR"}
		env := Environment{"TEST_VAR": EnvValue{Value: "success"}}

		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
	})

	t.Run("command with exit code", func(t *testing.T) {
		cmd := []string{"sh", "-c", "exit 42"}
		env := Environment{}

		code := RunCmd(cmd, env)
		if code != 42 {
			t.Errorf("Expected exit code 42, got %d", code)
		}
	})

	t.Run("nonexistent command", func(t *testing.T) {
		cmd := []string{"nonexistent-command-that-should-not-exist"}
		env := Environment{}

		code := RunCmd(cmd, env)
		if code != 1 {
			t.Errorf("Expected exit code 1 for nonexistent command, got %d", code)
		}
	})

	t.Run("no command specified", func(t *testing.T) {
		cmd := []string{}
		env := Environment{}

		code := RunCmd(cmd, env)
		if code != 1 {
			t.Errorf("Expected exit code 1 for empty command, got %d", code)
		}
	})
}
