package osutil

import (
	"os"
	"path/filepath"
)

// FilesRecursive returns all files in a given path and its subpaths.
func FilesRecursive(path string) (files []string, err error) {
	var fs []string

	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		fs = append(fs, path)

		return err
	})
	if err != nil {
		return nil, err
	}

	return fs, nil
}
