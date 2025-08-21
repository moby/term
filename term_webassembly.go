//go:build js

package term

import (
	"errors"
	"io"
	"os"
	"syscall/js"
)

var errUnsupported = errors.New("terminal operations are not supported on js")

// terminalState holds the platform-specific state / console mode for the terminal.
type terminalState struct{}

func stdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

func getFdInfo(in any) (uintptr, bool) {
	var inFd uintptr
	var isTerminalIn bool
	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminalIn = isTerminal(inFd)
	}
	return inFd, isTerminalIn
}

func getWinsize(fd uintptr) (*Winsize, error) {
	window := js.Global().Get("window")
	width := uint16(window.Get("innerWidth").Int())
	height := uint16(window.Get("innerHeight").Int())
	return &Winsize{Width: width, Height: height}, nil
}

func setWinsize(fd uintptr, ws *Winsize) error {
	window := js.Global().Get("window")
	window.Set("innerWidth", ws.Width)
	window.Set("innerHeight", ws.Height)

	return nil
}

func isTerminal(fd uintptr) bool {
	return true
}

func restoreTerminal(fd uintptr, state *State) error {
	if state == nil {
		return errors.New("invalid terminal state")
	}
	return nil
}

func saveState(fd uintptr) (*State, error) {
	return &State{}, nil
}

func disableEcho(fd uintptr, state *State) error {
	if state == nil {
		return errors.New("invalid terminal state")
	}
	return nil
}

func setRawTerminal(fd uintptr) (*State, error) {
	return &State{}, nil
}

func setRawTerminalOutput(fd uintptr) (*State, error) {
	return nil, errUnsupported
}
