//go:build windows
// +build windows

package windowsconsole

import (
	"io"
	"os"

	"github.com/Azure/go-ansiterm"
	"github.com/Azure/go-ansiterm/winterm"
)

// ansiWriter wraps a standard output file (e.g., os.Stdout) providing ANSI sequence translation.
type ansiWriter struct {
	file           *os.File
	fd             uintptr
	command        []byte
	escapeSequence []byte
	parser         *ansiterm.AnsiParser
}

// NewAnsiWriter returns an io.Writer that provides VT100 terminal emulation on top of a
// Windows console output handle.
func NewAnsiWriter(nFile int) io.Writer {
	file, fd := winterm.GetStdFile(nFile)
	handler := winterm.CreateWinEventHandler(fd, file)
	if handler == nil {
		return nil
	}

	return &ansiWriter{
		file:           file,
		fd:             fd,
		command:        make([]byte, 0, ansiterm.ANSI_MAX_CMD_LENGTH),
		escapeSequence: []byte(ansiterm.KEY_ESC_CSI),
		parser:         ansiterm.CreateParser("Ground", handler),
	}
}

func (aw *ansiWriter) Fd() uintptr {
	return aw.fd
}

// Write writes len(p) bytes from p to the underlying data stream.
func (aw *ansiWriter) Write(p []byte) (total int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	return aw.parser.Parse(p)
}
