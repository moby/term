//go:build js

package term

import (
	"errors"
	"io"
	"os"
)

var errUnsupported = errors.New("terminal operations are not supported on js")

// terminalState holds the platform-specific state / console mode for the terminal.
type terminalState struct{}

func stdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

func getFdInfo(in any) (uintptr, bool) {
	if file, ok := in.(*os.File); ok {
		return file.Fd(), false
	}
	return 0, false
}

func getWinsize(fd uintptr) (*Winsize, error) {
	return nil, errUnsupported
}

func setWinsize(fd uintptr, ws *Winsize) error {
	return errUnsupported
}

func isTerminal(fd uintptr) bool {
	return false
}

func restoreTerminal(fd uintptr, state *State) error {
	return errUnsupported
}

func saveState(fd uintptr) (*State, error) {
	return nil, errUnsupported
}

func disableEcho(fd uintptr, state *State) error {
	return errUnsupported
}

func setRawTerminal(fd uintptr) (*State, error) {
	return nil, errUnsupported
}

func setRawTerminalOutput(uintptr) (*State, error) {
	return nil, errUnsupported
}
