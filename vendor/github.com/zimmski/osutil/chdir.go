package osutil

import (
	"os"
)

// Chdir temporarily changes to the given working directory while calling the given function.
func Chdir(workingDirectory string, call func() error) (err error) {
	var owd string

	owd, err = os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		e := os.Chdir(owd)
		if err == nil {
			err = e
		}
	}()
	err = os.Chdir(workingDirectory)
	if err != nil {
		return err
	}

	return call()
}
