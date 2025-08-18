package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if err := Copy(from, to, offset, limit); err != nil {
		switch {
		case errors.Is(err, ErrIllegalArgument):
			fmt.Printf("Invalid arguments: %v\n", err)
		case errors.Is(err, ErrUnsupportedFile):
			fmt.Printf("Unsupported file: %v\n", err)
		case errors.Is(err, ErrOffsetExceedsFileSize):
			fmt.Printf("Offset error: %v\n", err)
		default:
			fmt.Printf("Copy failed: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("Copy completed successfully")
}
