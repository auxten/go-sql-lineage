// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package utils diff implements a Diff function that compare two inputs
// using the 'diff' tool.
package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

// Diff returns diff of two arrays of bytes in diff tool format.
func Diff(prefix string, b1, b2 []byte) ([]byte, error) {
	f1, err := writeTempFile(prefix, b1)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile(prefix, b2)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2)

	cmd := "diff"
	if runtime.GOOS == "plan9" {
		cmd = "/bin/ape/diff"
	}

	data, err := exec.Command(cmd, f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return data, err
}

func writeTempFile(prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	if _, err = file.Write(data); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	if _, err = file.Write([]byte{'\n'}); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}
