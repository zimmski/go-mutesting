package osutil

import (
	"syscall"
)

// SetRLimitFiles temporarily changes the file descriptor resource limit while calling the given function.
func SetRLimitFiles(limit uint64, call func(limit uint64)) (err error) {
	var tmp syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &tmp); err != nil {
		return nil
	}
	defer func() {
		if err == nil {
			err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &tmp)
		}
	}()

	if err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{
		Cur: limit,
		Max: tmp.Max,
	}); err != nil {
		return err
	}

	call(limit)

	return nil
}
