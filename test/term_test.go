//go:build !windows
// +build !windows

package test

import (
	"os"
	"reflect"
	"testing"

	cpty "github.com/creack/pty"
	"github.com/moby/term"
)

func newTTYForTest(t *testing.T) *os.File {
	t.Helper()
	pty, tty, err := cpty.Open()
	if err != nil {
		t.Fatalf("error creating pty: %v", err)
	} else {
		t.Cleanup(func() {
			_ = pty.Close()
			_ = tty.Close()
		})
	}

	return pty
}

func newTempFile(t *testing.T) *os.File {
	t.Helper()
	tmpFile, err := os.CreateTemp(t.TempDir(), "temp")
	if err != nil {
		t.Fatalf("error creating tempfile: %v", err)
	} else {
		t.Cleanup(func() { _ = tmpFile.Close() })
	}
	return tmpFile
}

func TestGetWinsize(t *testing.T) {
	tty := newTTYForTest(t)
	winSize, err := term.GetWinsize(tty.Fd())
	if err != nil {
		t.Error(err)
	}
	if winSize == nil {
		t.Fatal("winSize is nil")
	}

	newSize := term.Winsize{Width: 200, Height: 200}
	err = term.SetWinsize(tty.Fd(), &newSize)
	if err != nil {
		t.Fatal(err)
	}
	winSize, err = term.GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*winSize, newSize) {
		t.Fatalf("expected: %+v, got: %+v", newSize, *winSize)
	}
}

func TestSetWinsize(t *testing.T) {
	tty := newTTYForTest(t)
	winSize, err := term.GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if winSize == nil {
		t.Fatal("winSize is nil")
	}
	newSize := term.Winsize{Width: 200, Height: 200}
	err = term.SetWinsize(tty.Fd(), &newSize)
	if err != nil {
		t.Fatal(err)
	}
	winSize, err = term.GetWinsize(tty.Fd())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*winSize, newSize) {
		t.Fatalf("expected: %+v, got: %+v", newSize, *winSize)
	}
}

func TestGetFdInfo(t *testing.T) {
	tty := newTTYForTest(t)
	inFd, isTerm := term.GetFdInfo(tty)
	if inFd != tty.Fd() {
		t.Errorf("expected: %d, got: %d", tty.Fd(), inFd)
	}
	if !isTerm {
		t.Error("expected file-descriptor to be a terminal")
	}
	tmpFile := newTempFile(t)
	inFd, isTerm = term.GetFdInfo(tmpFile)
	if inFd != tmpFile.Fd() {
		t.Errorf("expected: %d, got: %d", tty.Fd(), inFd)
	}
	if isTerm {
		t.Error("expected file-descriptor to not be a terminal")
	}
}

func TestIsTerminal(t *testing.T) {
	tty := newTTYForTest(t)
	isTerm := term.IsTerminal(tty.Fd())
	if !isTerm {
		t.Error("expected file-descriptor to be a terminal")
	}
	tmpFile := newTempFile(t)
	isTerm = term.IsTerminal(tmpFile.Fd())
	if isTerm {
		t.Error("expected file-descriptor to not be a terminal")
	}
}

func TestSaveState(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := term.SaveState(tty.Fd())
	if err != nil {
		t.Error(err)
	}
	if state == nil {
		t.Fatal("state is nil")
	}
	tty = newTTYForTest(t)
	err = term.RestoreTerminal(tty.Fd(), state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDisableEcho(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := term.SetRawTerminal(tty.Fd())
	defer func() {
		if err := term.RestoreTerminal(tty.Fd(), state); err != nil {
			t.Error(err)
		}
	}()
	if err != nil {
		t.Error(err)
	}
	if state == nil {
		t.Fatal("state is nil")
	}
	err = term.DisableEcho(tty.Fd(), state)
	if err != nil {
		t.Fatal(err)
	}
}
