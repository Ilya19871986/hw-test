package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidEnvName = errors.New("invalid environment variable name: contains '='")

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	env := make(Environment)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			return nil, ErrInvalidEnvName
		}

		path := filepath.Join(dir, name)
		value, err := readEnvFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", name, err)
		}

		env[name] = value
	}

	return env, nil
}

// readEnvFile reads an environment file and returns EnvValue.
func readEnvFile(path string) (EnvValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return EnvValue{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return EnvValue{}, fmt.Errorf("failed to get file info: %w", err)
	}

	if info.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return EnvValue{}, fmt.Errorf("failed to read file content: %w", err)
	}

	// Удаляем ВСЕ пробельные символы в конце строки.
	line = strings.TrimRight(line, " \t\n\r")

	// Заменяем нулевые байты.
	line = strings.ReplaceAll(line, "\x00", "\n")

	return EnvValue{Value: line}, nil
}

// PrepareEnv prepares environment variables for command execution.
func (e Environment) PrepareEnv() []string {
	var env []string
	currentEnv := os.Environ()

	// Сначала добавляем существующие переменные, кроме тех, которые будут перезаписаны или удалены.
	for _, envVar := range currentEnv {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) != 2 {
			continue
		}
		name := pair[0]
		if _, exists := e[name]; !exists {
			env = append(env, envVar)
		}
	}

	// Затем добавляем новые переменные.
	for name, envValue := range e {
		if !envValue.NeedRemove {
			env = append(env, name+"="+envValue.Value)
		}
		// Для NeedRemove просто не добавляем переменную (уже исключили из currentEnv).
	}

	return env
}
