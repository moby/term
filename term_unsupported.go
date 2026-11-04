//go:build js || wasip1

package term

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

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
	return nil, fmt.Errorf("terminal: GetWinsize not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func setWinsize(fd uintptr, ws *Winsize) error {
	return fmt.Errorf("terminal: SetWinsize not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func isTerminal(fd uintptr) bool {
	return false
}

func restoreTerminal(fd uintptr, state *State) error {
	return fmt.Errorf("terminal: RestoreTerminal not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func saveState(fd uintptr) (*State, error) {
	return nil, fmt.Errorf("terminal: SaveState not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func disableEcho(fd uintptr, state *State) error {
	return fmt.Errorf("terminal: DisableEcho not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func setRawTerminal(fd uintptr) (*State, error) {
	return fmt.Errorf("terminal: SetRawTerminal not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func setRawTerminalOutput(uintptr) (*State, error) {
	return fmt.Errorf("terminal: SetRawTerminalOutput not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
