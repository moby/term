package term_test

import (
	"testing"

	"github.com/moby/term"
)

type testTerminalInfo struct {
	fd         uintptr
	isTerminal bool
}

func (t testTerminalInfo) FD() uintptr {
	return t.fd
}

func (t testTerminalInfo) IsTerminal() bool {
	return t.isTerminal
}

func TestGetFdInfoStream(t *testing.T) {
	for _, tc := range []struct {
		name       string
		isTerminal bool
	}{
		{name: "terminal", isTerminal: true},
		{name: "non-terminal", isTerminal: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Stream wrappers, such as docker/cli's streams.Out, provide
			// their own terminal status, regardless of the underlying
			// descriptor.
			stream := testTerminalInfo{
				fd:         123,
				isTerminal: tc.isTerminal,
			}

			fd, isTerm := term.GetFdInfo(stream)
			if fd != stream.FD() || isTerm != stream.IsTerminal() {
				t.Errorf("GetFdInfo() = (%d, %t), want (%d, %t)",
					fd, isTerm, stream.FD(), stream.IsTerminal())
			}
		})
	}
}
