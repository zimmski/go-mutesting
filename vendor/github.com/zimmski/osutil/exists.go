package osutil

import (
	"errors"
	"os"
)

var (
	// ErrNotADirectory indicates that the given directory does not exist.
	ErrNotADirectory = errors.New("not a directory")
	// ErrNotAFile indicates thate the given file does not exist.
	ErrNotAFile = errors.New("not a file")
)

// Stat retuns a FileInfo structure describing the given file.
func Stat(filepath string) (os.FileInfo, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return fi, nil
}

// DirExists checks if a directory exists.
func DirExists(path string) error {
	fi, err := Stat(path)
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return ErrNotADirectory
	}

	return nil
}

// FileExists checks if a file exists.
func FileExists(filepath string) error {
	fi, err := Stat(filepath)
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return ErrNotAFile
	}

	return nil
}
