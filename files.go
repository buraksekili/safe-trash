package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

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
func move(destination string, filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return fmt.Errorf("The input is a directory. Directory operation doesn't allowed for your safeness.\n")
	case mode.IsRegular():
		_, fname := filepath.Split(filename)
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("couldn't open '%s' file: %s", filename, err)
		}
		defer file.Close()

		destFile, err := os.Create(filepath.Join(destination, fname))
		if err != nil {
			return fmt.Errorf("couldn't open destination file: %s\n", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			return fmt.Errorf("error in copy operation: %s", err)
		}
		return os.Remove(filename)
	}
	return nil
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
