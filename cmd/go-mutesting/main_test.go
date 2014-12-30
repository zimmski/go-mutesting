package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	saveStderr := os.Stderr
	saveStdout := os.Stdout
	saveCwd, err := os.Getwd()
	assert.Nil(t, err)

	r, w, err := os.Pipe()
	assert.Nil(t, err)

	os.Stderr = w
	os.Stdout = w
	os.Chdir("../../example")

	bufChannel := make(chan string)

	go func() {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, r)
		r.Close()
		assert.Nil(t, err)

		bufChannel <- buf.String()
	}()

	exitCode := mainCmd([]string{"--exec", "../scripts/simple.sh", "--exec-timeout", "1", "./..."})

	w.Close()

	os.Stderr = saveStderr
	os.Stdout = saveStdout
	os.Chdir(saveCwd)

	out := <-bufChannel

	assert.Equal(t, returnOk, exitCode)
	assert.Contains(t, out, "T1he mutation score is 0.750000 (6 passed, 2 failed, 0 skipped, total is 8)")
}
