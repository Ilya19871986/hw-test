package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	err := Copy(from, to, offset, limit)
	if err != nil {
		switch err {
		case ErrIllegalArgument:
			panic(ErrIllegalArgument)
		case ErrUnsupportedFile:
			panic(ErrUnsupportedFile)
		case ErrOffsetExceedsFileSize:
			panic(ErrUnsupportedFile)
		default:
			panic(err)
		}
	}
	fmt.Println("Copy completed successfully")
}
