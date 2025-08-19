package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name     string
		content  string
		expected EnvValue
	}{
		{"FOO", "123", EnvValue{Value: "123"}},
		{"BAR", "value\n", EnvValue{Value: "value"}},
		{"EMPTY", "", EnvValue{NeedRemove: true}},
		{"WITH_SPACES", "  value  \t\n", EnvValue{Value: "  value"}},
		{"WITH_NULL", "first\x00second", EnvValue{Value: "first\nsecond"}},
		{"EMPTY_LINE", "\n", EnvValue{Value: ""}},
	}

	for _, tc := range testCases {
		path := filepath.Join(tmpDir, tc.name)
		if tc.content == "" {
			file, err := os.Create(path)
			if err != nil {
				t.Fatalf("Failed to create empty file: %v", err)
			}
			file.Close()
		} else {
			err := os.WriteFile(path, []byte(tc.content), 0o644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}
		}
	}

	env, err := ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	if len(env) != len(testCases) {
		t.Errorf("Expected %d env vars, got %d", len(testCases), len(env))
	}

	for _, tc := range testCases {
		val, ok := env[tc.name]
		if !ok {
			t.Errorf("Expected env var %q not found", tc.name)
			continue
		}

		if val != tc.expected {
			t.Errorf("For %q expected %v, got %v", tc.name, tc.expected, val)
		}
	}

	invalidDir := filepath.Join(tmpDir, "nonexistent")
	_, err = ReadDir(invalidDir)
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}

	invalidFile := filepath.Join(tmpDir, "INVALID=NAME")
	err = os.WriteFile(invalidFile, []byte("value"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	_, err = ReadDir(tmpDir)
	if err == nil {
		t.Error("Expected error for invalid environment variable name")
	} else if !errors.Is(err, ErrInvalidEnvName) {
		t.Errorf("Expected ErrInvalidEnvName, got %v", err)
	}
}

func TestPrepareEnv(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range originalEnv {
			pair := strings.SplitN(e, "=", 2)
			if len(pair) == 2 {
				os.Setenv(pair[0], pair[1])
			}
		}
	}()

	os.Clearenv()
	os.Setenv("EXISTING_VAR", "old_value")
	os.Setenv("TO_REMOVE", "should_be_removed")

	env := Environment{
		"FOO":          EnvValue{Value: "123"},
		"BAR":          EnvValue{Value: "value"},
		"TO_REMOVE":    EnvValue{NeedRemove: true},
		"EXISTING_VAR": EnvValue{Value: "new_value"},
	}

	envVars := env.PrepareEnv()

	expectedVars := map[string]string{
		"FOO":          "123",
		"BAR":          "value",
		"EXISTING_VAR": "new_value",
	}

	for name, expectedValue := range expectedVars {
		found := false
		for _, envVar := range envVars {
			if strings.HasPrefix(envVar, name+"=") {
				parts := strings.SplitN(envVar, "=", 2)
				if len(parts) == 2 && parts[1] == expectedValue {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("Expected env var %s=%s not found", name, expectedValue)
		}
	}

	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "TO_REMOVE=") {
			t.Errorf("Variable TO_REMOVE should be removed but found: %q", envVar)
		}
	}
}
