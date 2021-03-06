package main

import (
	"fmt"
	"os"
	"strings"
)

type Operation interface{}

type HelpOp struct{}

type FilenameOp struct {
	Names []string
}

type ListOp struct{}

type UnknownOp struct{}

func parseFlags() (Operation, error) {
	var flags []string = os.Args[1:]
	if len(flags) == 0 {
		return UnknownOp{}, fmt.Errorf("expects filename argument. check --help or -h")
	}

	var op = strings.TrimSpace(flags[0])

	// So, all flags must be file name.
	// Operation flags such as `-l` do not allowed between filenames.
	if !strings.HasPrefix(op, "-") && !strings.HasPrefix(op, "--") {
		var fnames []string
		for _, fn := range flags {
			if len(fn) > 0 {
				fnames = append(fnames, fn)
			}
		}
		return FilenameOp{Names: fnames}, nil
	}

	if op == "-h" || op == "--help" {
		return HelpOp{}, nil
	}

	if op == "-l" || op == "--list" {
		return ListOp{}, nil
	}

	return UnknownOp{}, fmt.Errorf("unknown operation flag\n")
}
