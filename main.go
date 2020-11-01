package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Operation interface{}

type HelpOp struct{}

type FilenameOp struct {
	Name string
}

type ListOp struct{}

type UnknownOp struct{}

func main() {

	op, err := parseFlags()
	if err != nil {
		fmt.Printf("Error while parsing flags: %s\n", err)
		return
	}

	switch v := op.(type) {
	case HelpOp:
		printHelp()

	case FilenameOp:
		trashPath, err := trashDir()
		if err != nil {
			fmt.Println("error occurred as ", err)
			return
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("error: ", err)
			return
		}

		err = move(cwd, trashPath, v.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully moved item to trash")

	case ListOp:
		if err := listCwd(); err != nil {
			fmt.Println(err)
		}

	case UnknownOp:
		fmt.Println("Unknown flag")
	}

}

// trashDir function creates a folder for the safe-trash if it
// doesn't exist. Returns path of the safe-trash folder and
// the error.
func trashDir() (string, error) {
	self, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unsupported OS: %s", err)
	}

	// path is for the trash folder.
	path := filepath.Join(self.HomeDir, ".safe-trash")
	if len(path) == 0 {
		return "", fmt.Errorf("unkown path as %s", path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("couldn't create trash folder: %s", err)
		}
		fmt.Println(".safe-trash folder is created at ", path)
		return path, nil
	}

	fmt.Printf("%s  will be used\n", path)
	return path, nil
}

// move function changes item location from source to destination.
// Returns error.
func move(source string, destination string, filename string) error {
	sourcePath := filepath.Join(source, filename)
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open '%s' file: %s", sourcePath, err)
	}

	// if destFile exists, keep all of the copies
	destFile, err := os.Create(filepath.Join(destination, filename))
	if err != nil {
		err = file.Close()
		return fmt.Errorf("couldn't open destination file: %s\n", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)

	_ = file.Close()
	if err != nil {
		return fmt.Errorf("error in copy operation: %s", err)
	}

	err = os.Remove(filepath.Join(source, filename))
	return nil
}

func parseFlags() (Operation, error) {
	var flags []string = os.Args[1:]
	if len(flags) == 0 {
		return UnknownOp{}, fmt.Errorf("expects filename argument. check --help or -h\n")
	}

	if len(flags) > 1 {
		return UnknownOp{}, fmt.Errorf("expects one argument at a time. check --help or -h\n")
	}

	var op = flags[0]
	if op == "-h" || op == "--help" {
		return HelpOp{}, nil
	}

	if op == "-l" || op == "--list" {
		return ListOp{}, nil
	}

	if !strings.HasPrefix(op, "-") {
		return FilenameOp{Name: op}, nil
	}

	return UnknownOp{}, fmt.Errorf("unknown operation flag\n")

}

func listCwd() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	return nil
}

func printHelp() {
	help := "USAGE:\n" +
		"\t-l, --list\t List the files under the current directory.\n" +
		"\t<FILE_NAME>\t The filename to be deleted.\n" +
		"\t-h, --help\t Displays this help message.\n"

	fmt.Printf("%s\n", help)
}
