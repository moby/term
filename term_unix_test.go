//go:build !windows

package term

import (
	"os"
	"reflect"
	"testing"

	cpty "github.com/creack/pty"
)

func newTTYForTest(t *testing.T) *os.File {
	t.Helper()
	pty, tty, err := cpty.Open()
	if err != nil {
		t.Fatalf("error creating pty: %v", err)
	}
	t.Cleanup(func() {
		_ = pty.Close()
		_ = tty.Close()
	})

	return pty
}

func newTempFile(t *testing.T) *os.File {
	t.Helper()
	tmpFile, err := os.CreateTemp(t.TempDir(), "temp")
	if err != nil {
		t.Fatalf("error creating tempfile: %v", err)
	}
	t.Cleanup(func() { _ = tmpFile.Close() })
	return tmpFile
}

func TestGetWinsize(t *testing.T) {
	tty := newTTYForTest(t)
	winSize, err := GetWinsize(tty.Fd())
	if err != nil {
		t.Error(err)
	}
	if winSize == nil {
		t.Fatal("winSize is nil")
	}

	newSize := Winsize{Width: 200, Height: 200}
	err = SetWinsize(tty.Fd(), &newSize)
	if err != nil {
		t.Fatal(err)
	}
	winSize, err = GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*winSize, newSize) {
		t.Fatalf("expected: %+v, got: %+v", newSize, *winSize)
	}
}

func TestSetWinsize(t *testing.T) {
	tty := newTTYForTest(t)
	winSize, err := GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if winSize == nil {
		t.Fatal("winSize is nil")
	}
	newSize := Winsize{Width: 200, Height: 200}
	err = SetWinsize(tty.Fd(), &newSize)
	if err != nil {
		t.Fatal(err)
	}
	winSize, err = GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*winSize, newSize) {
		t.Fatalf("expected: %+v, got: %+v", newSize, *winSize)
	}
}

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

func TestGetFdInfo(t *testing.T) {
	tty := newTTYForTest(t)
	tmpFile := newTempFile(t)

	for _, tc := range []struct {
		name       string
		stream     any
		wantFD     uintptr
		wantIsTerm bool
	}{
		{
			name:       "tty",
			stream:     tty,
			wantFD:     tty.Fd(),
			wantIsTerm: true,
		},
		{
			name:       "file",
			stream:     tmpFile,
			wantFD:     tmpFile.Fd(),
			wantIsTerm: false,
		},
		// Stream wrappers, such as docker/cli's streams.Out, provide their
		// own terminal status. Use a TTY descriptor to verify that the
		// reported status takes precedence over OS-level detection.
		{
			name:       "stream terminal",
			stream:     testTerminalInfo{fd: tty.Fd(), isTerminal: true},
			wantFD:     tty.Fd(),
			wantIsTerm: true,
		},
		{
			name:       "stream non-terminal",
			stream:     testTerminalInfo{fd: tty.Fd(), isTerminal: false},
			wantFD:     tty.Fd(),
			wantIsTerm: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fd, isTerm := GetFdInfo(tc.stream)
			if fd != tc.wantFD || isTerm != tc.wantIsTerm {
				t.Errorf("GetFdInfo() = (%d, %t), want (%d, %t)", fd, isTerm, tc.wantFD, tc.wantIsTerm)
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	tty := newTTYForTest(t)
	isTerm := IsTerminal(tty.Fd())
	if !isTerm {
		t.Error("expected file-descriptor to be a terminal")
	}
	tmpFile := newTempFile(t)
	isTerm = IsTerminal(tmpFile.Fd())
	if isTerm {
		t.Error("expected file-descriptor to not be a terminal")
	}
}

func TestSaveState(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := SaveState(tty.Fd())
	if err != nil {
		t.Error(err)
	}
	if state == nil {
		t.Fatal("state is nil")
	}
	tty = newTTYForTest(t)
	err = RestoreTerminal(tty.Fd(), state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDisableEcho(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := SetRawTerminal(tty.Fd())
	defer func() {
		if err := RestoreTerminal(tty.Fd(), state); err != nil {
			t.Error(err)
		}
	}()
	if err != nil {
		t.Error(err)
	}
	if state == nil {
		t.Fatal("state is nil")
	}
	err = DisableEcho(tty.Fd(), state)
	if err != nil {
		t.Fatal(err)
	}
}
