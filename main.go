package main

import (
	"log"
	"os"

	"github.com/fatih/color"

	"github.com/buraksekili/selog"
)

func main() {

	logger := log.New(os.Stdout, "safe-trash", log.LstdFlags)
	l := selog.NewLogger(logger)
	
	op, err := parseFlags()
	if err != nil {
		l.Fatal("couldn't parse flag: %v\n", err)
	}

	switch v := op.(type) {
	case HelpOp:
		printHelp()
	case FilenameOp:
		trashPath, err := trashDir()
		if err != nil {
			l.Fatal("couldn't generate trash directory: %v\n", err)
		}

		for _, fn := range v.Names {
			err = move(trashPath, fn)
			if err != nil {
				l.Fatal("couldn't change the cwd: %v\n", err)
			}
			color.Green("Successfully moved %s to trash", fn)
		}

	case ListOp:
		if err := listCwd(); err != nil {
			l.Fatal("couldn't list cwd", err)
		}
	case UnknownOp:
		l.Error("Unknown flag\n")
	}
}
