package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"os/user"
)

func main() {
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

	err = move(cwd, trashPath)
	if err != nil {
		fmt.Println(err)
		return
	}
}

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

func move(source string, destination string) error {
	fmt.Println()
	file, err := os.Open(filepath.Join(source, "ex.py"))
	if err != nil {
		return fmt.Errorf("couldn't open '%s' file: %s", source, err)
	}

	// if destFile exists, keep all of the copies
	destFile, err := os.Create(filepath.Join(destination, "ex.py"))
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

	err = os.Remove(filepath.Join(source, "ex.py"))

	return nil
}
