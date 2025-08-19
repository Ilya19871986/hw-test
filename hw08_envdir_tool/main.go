package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <envdir> <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	envDir := os.Args[1]
	commandArgs := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading environment: %v\n", err)
		os.Exit(1)
	}

	exitCode := RunCmd(commandArgs, env)
	os.Exit(exitCode)
}
