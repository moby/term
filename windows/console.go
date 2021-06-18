// +build windows

package windowsconsole

import (
	"os"

	"golang.org/x/sys/windows"
)

type fd interface {
	Fd() uintptr
}

// GetHandleInfo returns file descriptor and bool indicating whether the file is a console.
func GetHandleInfo(in interface{}) (uintptr, bool) {
	var inFd uintptr
	var isTerminal bool

	if file, ok := in.(fd); ok {
		inFd = file.Fd()
		isTerminal = isConsole(inFd)
	}
	return inFd, isTerminal
}

// IsConsole returns true if the given file descriptor is a Windows Console.
// The code assumes that GetConsoleMode will return an error for file descriptors that are not a console.
// Deprecated: use golang.org/x/sys/windows.GetConsoleMode() or golang.org/x/term.IsTerminal()
var IsConsole = isConsole

func isConsole(fd uintptr) bool {
	var mode uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &mode)
	return err == nil
}
