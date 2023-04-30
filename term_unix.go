//go:build !windows
// +build !windows

package term

import (
	"errors"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

// ErrInvalidState is returned if the state of the terminal is invalid.
//
// Deprecated: ErrInvalidState is no longer used.
var ErrInvalidState = errors.New("Invalid terminal state")

// State represents the state of the terminal.
type State struct {
	termios unix.Termios
}

// Winsize represents the size of the terminal window.
type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16
	y      uint16
}

// StdStreams returns the standard streams (stdin, stdout, stderr).
func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

// GetFdInfo returns the file descriptor for an os.File and indicates whether the file represents a terminal.
func GetFdInfo(in interface{}) (uintptr, bool) {
	var inFd uintptr
	var isTerminalIn bool
	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminalIn = IsTerminal(inFd)
	}
	return inFd, isTerminalIn
}

// GetWinsize returns the window size based on the specified file descriptor.
func GetWinsize(fd uintptr) (*Winsize, error) {
	uws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	ws := &Winsize{Height: uws.Row, Width: uws.Col, x: uws.Xpixel, y: uws.Ypixel}
	return ws, err
}

// SetWinsize tries to set the specified window size for the specified file descriptor.
func SetWinsize(fd uintptr, ws *Winsize) error {
	uws := &unix.Winsize{Row: ws.Height, Col: ws.Width, Xpixel: ws.x, Ypixel: ws.y}
	return unix.IoctlSetWinsize(int(fd), unix.TIOCSWINSZ, uws)
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool {
	_, err := tcget(fd)
	return err == nil
}

// RestoreTerminal restores the terminal connected to the given file descriptor
// to a previous state.
func RestoreTerminal(fd uintptr, state *State) error {
	if state == nil {
		return errors.New("invalid terminal state")
	}
	return tcset(fd, &state.termios)
}

// SaveState saves the state of the terminal connected to the given file descriptor.
func SaveState(fd uintptr) (*State, error) {
	termios, err := tcget(fd)
	if err != nil {
		return nil, err
	}
	return &State{termios: *termios}, nil
}

// DisableEcho applies the specified state to the terminal connected to the file
// descriptor, with echo disabled.
func DisableEcho(fd uintptr, state *State) error {
	newState := state.termios
	newState.Lflag &^= unix.ECHO

	if err := tcset(fd, &newState); err != nil {
		return err
	}
	return nil
}

// SetRawTerminal puts the terminal connected to the given file descriptor into
// raw mode and returns the previous state. On UNIX, this puts both the input
// and output into raw mode. On Windows, it only puts the input into raw mode.
func SetRawTerminal(fd uintptr) (*State, error) {
	oldState, err := MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return oldState, err
}

// SetRawTerminalOutput puts the output of terminal connected to the given file
// descriptor into raw mode. On UNIX, this does nothing and returns nil for the
// state. On Windows, it disables LF -> CRLF translation.
func SetRawTerminalOutput(fd uintptr) (*State, error) {
	return nil, nil
}

func tcget(fd uintptr) (*unix.Termios, error) {
	p, err := unix.IoctlGetTermios(int(fd), getTermios)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func tcset(fd uintptr, p *unix.Termios) error {
	return unix.IoctlSetTermios(int(fd), setTermios, p)
}
