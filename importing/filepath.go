package importing

/*

This file holds lots of code of the golint project https://github.com/golang/lint and some code of a pull request of mine https://github.com/golang/lint/pull/76
This is just temporary until I have time to clean up this code and make a more general solution for go-commands as I stated here https://github.com/kisielk/errcheck/issues/45#issuecomment-57732642

so TODO and FIXME. Heck I also give you a WORKAROUND.

*/

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// FilesOfArgs returns all available Go files given a list of packages, directories and files which can embed patterns.
func FilesOfArgs(args []string) []string {
	var filenames []string

	if len(args) == 0 {
		filenames = append(filenames, checkDir(".")...)
	} else {
		for _, arg := range args {
			if strings.HasSuffix(arg, "/...") && isDir(arg[:len(arg)-4]) {
				for _, dirname := range allPackagesInFS(arg) {
					filenames = append(filenames, checkDir(dirname)...)
				}
			} else if isDir(arg) {
				filenames = append(filenames, checkDir(arg)...)
			} else if exists(arg) {
				filenames = append(filenames, arg)
			} else {
				for _, pkgname := range importPaths([]string{arg}) {
					filenames = append(filenames, checkPackage(pkgname)...)
				}
			}
		}
	}

	fileLookup := make(map[string]struct{})
	pkgs := make(map[string]map[string]struct{})
	for _, filename := range filenames {
		if _, ok := fileLookup[filename]; ok {
			continue
		}

		if strings.HasSuffix(filename, "_test.go") { // ignore test files
			continue
		}

		if !exists(filename) {
			fmt.Printf("%q does not exist", filename)

			continue
		}
		fileLookup[filename] = struct{}{}

		pkgName := path.Dir(filename)

		pkg, ok := pkgs[pkgName]
		if !ok {
			pkg = make(map[string]struct{})

			pkgs[pkgName] = pkg
		}

		pkg[filename] = struct{}{}
	}

	pkgsNames := make([]string, 0, len(pkgs))
	for name := range pkgs {
		pkgsNames = append(pkgsNames, name)
	}
	sort.Strings(pkgsNames)

	var files []string

	for _, name := range pkgsNames {
		var filenames []string
		for name := range pkgs[name] {
			filenames = append(filenames, name)
		}
		sort.Strings(filenames)

		files = append(files, filenames...)
	}

	return files
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func checkDir(dirname string) []string {
	pkg, err := build.ImportDir(dirname, 0)

	return checkImportedPackage(pkg, err)
}

func checkPackage(pkgname string) []string {
	pkg, err := build.Import(pkgname, ".", 0)

	return checkImportedPackage(pkg, err)
}

func checkImportedPackage(pkg *build.Package, err error) []string {
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return []string{}
		}
		fmt.Fprintln(os.Stderr, err)

		return []string{}
	}

	var files []string

	files = append(files, pkg.GoFiles...)

	joinDirWithFilenames(pkg.Dir, files)

	return files
}

func joinDirWithFilenames(dir string, files []string) {
	if dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(dir, f)
		}
	}
}
