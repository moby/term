//go:build windows
// +build windows

package windowsconsole

import (
	"os"

	"golang.org/x/sys/windows"
)

// GetHandleInfo returns file descriptor and bool indicating whether the file is a console.
func GetHandleInfo(in interface{}) (uintptr, bool) {
	switch t := in.(type) {
	case *ansiReader:
		return t.Fd(), true
	case *ansiWriter:
		return t.Fd(), true
	}

	var inFd uintptr
	var isTerminal bool

	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminal = isConsole(inFd)
	}
	return inFd, isTerminal
}

func isConsole(fd uintptr) bool {
	var mode uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &mode)
	return err == nil
}
