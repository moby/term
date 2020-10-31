//go:build !windows
// +build !windows

package term

import (
	"os"
	"testing"

	cpty "github.com/creack/pty"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
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
	return tty
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
	winSize, err := GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, winSize != nil)

	newSize := Winsize{Width: 200, Height: 200, x: winSize.x, y: winSize.y}
	err = SetWinsize(tty.Fd(), &newSize)
	assert.NilError(t, err)
	winSize, err = GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.DeepEqual(t, *winSize, newSize, cmpWinsize)
}

var cmpWinsize = cmp.AllowUnexported(Winsize{})

func TestSetWinsize(t *testing.T) {
	tty := newTTYForTest(t)
	winSize, err := GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, winSize != nil)
	newSize := Winsize{Width: 200, Height: 200, x: winSize.x, y: winSize.y}
	err = SetWinsize(tty.Fd(), &newSize)
	assert.NilError(t, err)
	winSize, err = GetWinsize(tty.Fd())
	assert.NilError(t, err)
	assert.DeepEqual(t, *winSize, newSize, cmpWinsize)
}

func TestGetFdInfo(t *testing.T) {
	tty := newTTYForTest(t)
	inFd, isTerminal := GetFdInfo(tty)
	assert.Equal(t, inFd, tty.Fd())
	assert.Equal(t, isTerminal, true)
	tmpFile := newTempFile(t)
	inFd, isTerminal = GetFdInfo(tmpFile)
	assert.Equal(t, inFd, tmpFile.Fd())
	assert.Equal(t, isTerminal, false)
}

func TestIsTerminal(t *testing.T) {
	tty := newTTYForTest(t)
	isTerminal := IsTerminal(tty.Fd())
	assert.Equal(t, isTerminal, true)
	tmpFile := newTempFile(t)
	isTerminal = IsTerminal(tmpFile.Fd())
	assert.Equal(t, isTerminal, false)
}

func TestSaveState(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := SaveState(tty.Fd())
	assert.NilError(t, err)
	assert.Assert(t, state != nil)
	tty = newTTYForTest(t)
	defer tty.Close()
	err = RestoreTerminal(tty.Fd(), state)
	assert.NilError(t, err)
}

func TestDisableEcho(t *testing.T) {
	tty := newTTYForTest(t)
	state, err := SetRawTerminal(tty.Fd())
	defer RestoreTerminal(tty.Fd(), state)
	assert.NilError(t, err)
	assert.Assert(t, state != nil)
	err = DisableEcho(tty.Fd(), state)
	assert.NilError(t, err)
}
