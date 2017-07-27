// Based on ssh/terminal:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows,!appengine

package logrus

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/windows"
)

var kernel32 = windows.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

const (
	enableProcessedOutput           = 0x0001
	enableWrapAtEolOutput           = 0x0002
	enableVirtualTerminalProcessing = 0x0004
)

func getVersion() (float64, error) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd := exec.Command("cmd", "ver")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return -1, err
	}

	// The output should be like "Microsoft Windows [Version XX.X.XXXXXX]"
	version := strings.Replace(stdout.String(), "\n", "", -1)
	version = strings.Replace(version, "\r\n", "", -1)

	x1 := strings.Index(version, "[Version")

	if x1 == -1 || strings.Index(version, "]") == -1 {
		return -1, errors.New("Can't determine Windows version")
	}

	return strconv.ParseFloat(version[x1+9:x1+13], 64)
}

func init() {
	ver, err := getVersion()
	if err != nil {
		return
	}

	// Activate Virtual Processing for Windows CMD
	// Info: https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	if ver >= 10 {
		handle := windows.Handle(os.Stderr.Fd())
		procSetConsoleMode.Call(uintptr(handle), enableProcessedOutput|enableWrapAtEolOutput|enableVirtualTerminalProcessing)
	}
}

// IsTerminal returns true if stderr's file descriptor is a terminal.
func IsTerminal(f io.Writer) bool {
	switch v := f.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
